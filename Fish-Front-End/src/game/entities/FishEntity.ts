import {
  Scene,
  Mesh,
  MeshBuilder,
  StandardMaterial,
  Color3,
  Vector3,
  TransformNode,
  SceneLoader,
} from '@babylonjs/core'
import type { Fish } from '../../types'

export class FishEntity {
  private scene: Scene
  public readonly fishData: Fish
  private root: TransformNode
  private healthBarFill: Mesh | null = null
  private healthBarBg: Mesh | null = null

  public currentHealth: number
  public isDead = false
  public id: number

  // Swimming path parameters
  private time = 0
  private startX: number
  private startZ: number
  private amplitude: number
  private frequency: number
  private direction: 1 | -1
  private yPos: number

  onDeath?: (fish: FishEntity) => void

  constructor(scene: Scene, fishData: Fish, spawnIndex: number) {
    this.scene = scene
    this.fishData = fishData
    this.currentHealth = fishData.health
    this.id = fishData.id

    this.root = new TransformNode(`fish-root-${fishData.id}-${spawnIndex}`, scene)

    // Randomize swimming path
    this.direction = Math.random() > 0.5 ? 1 : -1
    this.startX = this.direction === 1 ? -30 : 30
    this.startZ = (Math.random() - 0.5) * 20
    this.amplitude = 1 + Math.random() * 2
    this.frequency = 0.5 + Math.random() * 1.5
    this.yPos = -1 + (Math.random() - 0.5) * 2
    this.time = Math.random() * Math.PI * 2

    this.root.position = new Vector3(this.startX, this.yPos, this.startZ)

    this.loadMesh()
    this.createHealthBar()
  }

  private async loadMesh() {
    if (this.fishData.asset_path) {
      try {
        const result = await SceneLoader.ImportMeshAsync('', '', this.fishData.asset_path, this.scene)
        if (result.meshes.length > 0 && !this.isDead) {
          const importedRoot = result.meshes[0]
          importedRoot.parent = this.root
          importedRoot.scaling = new Vector3(0.5, 0.5, 0.5)
        }
      } catch {
        this.createFallbackMesh()
      }
    } else {
      this.createFallbackMesh()
    }
  }

  private createFallbackMesh() {
    if (this.isDead) return

    const bodyColors: Color3[] = [
      new Color3(0.2, 0.8, 0.9),
      new Color3(0.9, 0.5, 0.1),
      new Color3(0.8, 0.2, 0.6),
      new Color3(0.3, 0.9, 0.4),
      new Color3(0.9, 0.8, 0.1),
    ]
    const color = bodyColors[this.fishData.id % bodyColors.length]

    // Body
    const body = MeshBuilder.CreateSphere(
      `fish-body-${this.fishData.id}`,
      { diameterX: 1.2, diameterY: 0.6, diameterZ: 0.8 },
      this.scene,
    )
    const mat = new StandardMaterial(`fish-mat-${this.fishData.id}`, this.scene)
    mat.diffuseColor = color
    mat.specularColor = new Color3(0.5, 0.5, 0.5)
    body.material = mat
    body.parent = this.root

    // Tail fin
    const tail = MeshBuilder.CreateCylinder(
      `fish-tail-${this.fishData.id}`,
      { height: 0.1, diameterTop: 0, diameterBottom: 0.8, tessellation: 4 },
      this.scene,
    )
    tail.rotation.z = Math.PI / 2
    tail.position.x = -0.8
    tail.material = mat
    tail.parent = this.root
  }

  private createHealthBar() {
    const barWidth = 1.5
    const barHeight = 0.15
    const yOffset = 1.2

    // Background bar (red)
    this.healthBarBg = MeshBuilder.CreatePlane(
      `hp-bg-${this.fishData.id}`,
      { width: barWidth, height: barHeight },
      this.scene,
    )
    const bgMat = new StandardMaterial(`hp-bg-mat-${this.fishData.id}`, this.scene)
    bgMat.diffuseColor = new Color3(0.6, 0.1, 0.1)
    bgMat.emissiveColor = new Color3(0.4, 0.05, 0.05)
    bgMat.disableLighting = true
    this.healthBarBg.material = bgMat
    this.healthBarBg.parent = this.root
    this.healthBarBg.position.y = yOffset
    this.healthBarBg.billboardMode = Mesh.BILLBOARDMODE_ALL

    // Filled bar (green)
    this.healthBarFill = MeshBuilder.CreatePlane(
      `hp-fill-${this.fishData.id}`,
      { width: barWidth, height: barHeight },
      this.scene,
    )
    const fillMat = new StandardMaterial(`hp-fill-mat-${this.fishData.id}`, this.scene)
    fillMat.diffuseColor = new Color3(0.1, 0.9, 0.3)
    fillMat.emissiveColor = new Color3(0.05, 0.6, 0.15)
    fillMat.disableLighting = true
    this.healthBarFill.material = fillMat
    this.healthBarFill.parent = this.root
    this.healthBarFill.position.y = yOffset
    this.healthBarFill.position.z = -0.01
    this.healthBarFill.billboardMode = Mesh.BILLBOARDMODE_ALL
  }

  private updateHealthBar() {
    if (!this.healthBarFill || !this.healthBarBg) return
    const ratio = Math.max(0, this.currentHealth / this.fishData.health)
    this.healthBarFill.scaling.x = ratio
    this.healthBarFill.position.x = (-(1 - ratio) * 1.5) / 2

    // Color shift from green to red
    const mat = this.healthBarFill.material as StandardMaterial
    mat.diffuseColor = new Color3(1 - ratio, ratio, 0.1)
    mat.emissiveColor = new Color3((1 - ratio) * 0.6, ratio * 0.6, 0.05)
  }

  takeDamage(amount: number): boolean {
    if (this.isDead) return false
    this.currentHealth -= amount
    this.updateHealthBar()

    if (this.currentHealth <= 0) {
      this.isDead = true
      this.onDeath?.(this)
      return true
    }
    return false
  }

  update(deltaTime: number) {
    if (this.isDead) return

    this.time += deltaTime * this.fishData.speed * 0.5

    const speed = this.fishData.speed * 2.5
    const x = this.startX + this.direction * this.time * speed
    const z = this.startZ + Math.sin(this.time * this.frequency) * this.amplitude
    const y = this.yPos + Math.sin(this.time * 0.3) * 0.3

    this.root.position = new Vector3(x, y, z)

    // Rotate to face movement direction with slight sine wobble
    const wobble = Math.sin(this.time * 2) * 0.1
    this.root.rotation.y = this.direction === 1 ? wobble : Math.PI + wobble

    // Reset when off screen
    if ((this.direction === 1 && x > 35) || (this.direction === -1 && x < -35)) {
      this.startX = this.direction === 1 ? -30 : 30
      this.startZ = (Math.random() - 0.5) * 20
      this.time = 0
    }
  }

  getRootMesh(): TransformNode {
    return this.root
  }

  getWorldPosition(): Vector3 {
    return this.root.getAbsolutePosition()
  }

  dispose() {
    this.root.getChildMeshes(false).forEach((m) => m.dispose())
    this.healthBarBg?.dispose()
    this.healthBarFill?.dispose()
    this.root.dispose()
  }
}
