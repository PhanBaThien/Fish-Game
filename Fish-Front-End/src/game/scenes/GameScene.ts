import {
  Engine,
  Scene,
  ArcRotateCamera,
  HemisphericLight,
  PointLight,
  Color3,
  Color4,
  Vector3,
  MeshBuilder,
  StandardMaterial,
  Texture,
  PickingInfo,
  ParticleSystem,
  SphereParticleEmitter,
} from '@babylonjs/core'
import { FishEntity } from '../entities/FishEntity'
import { BulletEntity } from '../entities/BulletEntity'
import type { Fish } from '../../types'

export interface GameSceneOptions {
  canvas: HTMLCanvasElement
  fishList: Fish[]
  onFishKilled?: (rewardMultiplier: number) => void
  onScore?: (points: number) => void
}

export class GameScene {
  private engine: Engine
  private scene: Scene
  private fishEntities: FishEntity[] = []
  private bullets: BulletEntity[] = []
  private options: GameSceneOptions
  private lastTime = 0

  constructor(options: GameSceneOptions) {
    this.options = options
    this.engine = new Engine(options.canvas, true, { preserveDrawingBuffer: true })
    this.scene = new Scene(this.engine)

    this.setupScene()
    this.spawnFish()
    this.setupInput()
    this.startRenderLoop()

    window.addEventListener('resize', this.handleResize)
  }

  private setupScene() {
    const scene = this.scene

    // Underwater dark blue background
    scene.clearColor = new Color4(0.02, 0.05, 0.15, 1)

    // Fog for underwater depth effect
    scene.fogMode = Scene.FOGMODE_EXP2
    scene.fogColor = new Color3(0.02, 0.05, 0.15)
    scene.fogDensity = 0.03

    // Camera: ArcRotate looking down at game area
    const camera = new ArcRotateCamera('cam', -Math.PI / 2, Math.PI / 3.5, 28, Vector3.Zero(), scene)
    camera.lowerRadiusLimit = 15
    camera.upperRadiusLimit = 45
    camera.lowerBetaLimit = 0.3
    camera.upperBetaLimit = Math.PI / 2.2
    camera.attachControl(this.options.canvas, false)

    // Ambient hemispheric light (underwater blue tint)
    const hemi = new HemisphericLight('hemi', new Vector3(0, 1, 0), scene)
    hemi.intensity = 0.6
    hemi.diffuse = new Color3(0.3, 0.7, 1)
    hemi.groundColor = new Color3(0.05, 0.1, 0.3)
    hemi.specular = new Color3(0.1, 0.2, 0.4)

    // Point lights for underwater caustics feel
    const light1 = new PointLight('light1', new Vector3(5, 8, 0), scene)
    light1.diffuse = new Color3(0.2, 0.8, 1)
    light1.intensity = 0.8
    light1.range = 25

    const light2 = new PointLight('light2', new Vector3(-5, 6, -5), scene)
    light2.diffuse = new Color3(0.1, 0.4, 0.8)
    light2.intensity = 0.5
    light2.range = 20

    // Ocean floor
    const ground = MeshBuilder.CreateGround('ground', { width: 80, height: 50, subdivisions: 10 }, scene)
    const groundMat = new StandardMaterial('groundMat', scene)
    groundMat.diffuseColor = new Color3(0.05, 0.15, 0.3)
    groundMat.specularColor = new Color3(0.1, 0.3, 0.5)
    ground.material = groundMat
    ground.position.y = -5

    // Invisible ground plane for ray picking bullet targets
    const shootPlane = MeshBuilder.CreatePlane('shootPlane', { width: 80, height: 50 }, scene)
    shootPlane.rotation.x = Math.PI / 2
    shootPlane.position.y = 0
    shootPlane.isPickable = true
    shootPlane.visibility = 0

    // Ambient bubble particles
    this.createBubbleParticles()
  }

  private createBubbleParticles() {
    const ps = new ParticleSystem('bubbles', 100, this.scene)
    ps.particleTexture = new Texture(
      'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==',
      this.scene,
    )
    ps.emitter = Vector3.Zero()
    const emitter = new SphereParticleEmitter(15)
    ps.particleEmitterType = emitter

    ps.color1 = new Color4(0.3, 0.7, 1, 0.6)
    ps.color2 = new Color4(0.1, 0.4, 0.8, 0.3)
    ps.colorDead = new Color4(0, 0.2, 0.5, 0)

    ps.minSize = 0.05
    ps.maxSize = 0.2
    ps.minLifeTime = 3
    ps.maxLifeTime = 8
    ps.emitRate = 10
    ps.gravity = new Vector3(0, 0.5, 0)
    ps.minEmitPower = 0.1
    ps.maxEmitPower = 0.5

    ps.start()
  }

  private spawnFish() {
    const { fishList } = this.options
    if (!fishList || fishList.length === 0) return

    // Spawn multiple instances of each fish type
    const totalFish = Math.max(fishList.length * 2, 8)
    for (let i = 0; i < totalFish; i++) {
      const fishData = fishList[i % fishList.length]
      this.addFish(fishData, i)
    }
  }

  addFish(fishData: Fish, spawnIndex = 0) {
    const entity = new FishEntity(this.scene, fishData, spawnIndex)
    entity.onDeath = (fish) => {
      this.handleFishDeath(fish)
    }
    this.fishEntities.push(entity)
  }

  private handleFishDeath(fish: FishEntity) {
    const pos = fish.getWorldPosition()
    this.spawnDeathParticles(pos)
    this.options.onFishKilled?.(fish.fishData.reward_multiplier)
    this.options.onScore?.(10)

    // Remove from array
    const idx = this.fishEntities.indexOf(fish)
    if (idx !== -1) {
      this.fishEntities.splice(idx, 1)
    }
    fish.dispose()

    // Respawn a random fish after a delay
    setTimeout(() => {
      const { fishList } = this.options
      if (fishList && fishList.length > 0 && !this.isDisposed) {
        const randomFish = fishList[Math.floor(Math.random() * fishList.length)]
        this.addFish(randomFish, Math.floor(Math.random() * 100))
      }
    }, 2000)
  }

  private spawnDeathParticles(position: Vector3) {
    const ps = new ParticleSystem(`death-${Date.now()}`, 50, this.scene)
    ps.particleTexture = new Texture(
      'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==',
      this.scene,
    )
    ps.emitter = position.clone()

    ps.color1 = new Color4(1, 0.8, 0.2, 1)
    ps.color2 = new Color4(1, 0.3, 0.1, 1)
    ps.colorDead = new Color4(0.5, 0.1, 0, 0)

    ps.minSize = 0.1
    ps.maxSize = 0.5
    ps.minLifeTime = 0.3
    ps.maxLifeTime = 0.8
    ps.emitRate = 200
    ps.gravity = new Vector3(0, -3, 0)
    ps.minEmitPower = 3
    ps.maxEmitPower = 8
    ps.updateSpeed = 0.02

    ps.start()
    setTimeout(() => {
      ps.stop()
      setTimeout(() => ps.dispose(), 2000)
    }, 300)
  }

  private setupInput() {
    this.scene.onPointerDown = (_evt, pickInfo) => {
      if (_evt.button !== 0) return
      this.shoot(pickInfo)
    }
  }

  private shoot(pickInfo: PickingInfo | null) {
    const camera = this.scene.activeCamera
    if (!camera) return

    let target: Vector3

    if (pickInfo?.hit && pickInfo.pickedPoint) {
      target = pickInfo.pickedPoint
    } else {
      // Fallback: shoot forward from camera
      const forward = camera.getDirection(Vector3.Forward())
      target = camera.position.add(forward.scale(20))
    }

    const origin = camera.position.clone()
    const bullet = new BulletEntity(this.scene, origin, target)
    this.bullets.push(bullet)
  }

  private checkBulletCollisions() {
    for (const bullet of this.bullets) {
      if (bullet.isDead) continue

      const bPos = bullet.getPosition()

      for (const fish of this.fishEntities) {
        if (fish.isDead) continue

        const fPos = fish.getWorldPosition()
        const dist = Vector3.Distance(bPos, fPos)

        if (dist < 1.2) {
          fish.takeDamage(10)
          bullet.destroy()
          break
        }
      }
    }

    // Clean up dead bullets
    this.bullets = this.bullets.filter((b) => !b.isDead)
  }

  private isDisposed = false

  private startRenderLoop() {
    this.lastTime = performance.now()

    this.engine.runRenderLoop(() => {
      const now = performance.now()
      const deltaTime = Math.min((now - this.lastTime) / 1000, 0.1)
      this.lastTime = now

      // Update fish
      for (const fish of this.fishEntities) {
        fish.update(deltaTime)
      }

      // Update bullets
      for (const bullet of this.bullets) {
        bullet.update(deltaTime)
      }

      this.checkBulletCollisions()

      this.scene.render()
    })
  }

  private handleResize = () => {
    this.engine.resize()
  }

  dispose() {
    this.isDisposed = true
    window.removeEventListener('resize', this.handleResize)

    for (const fish of this.fishEntities) {
      fish.dispose()
    }
    for (const bullet of this.bullets) {
      bullet.destroy()
    }

    this.scene.dispose()
    this.engine.dispose()
  }
}
