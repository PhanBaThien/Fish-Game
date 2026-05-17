import {
  Scene,
  Mesh,
  MeshBuilder,
  StandardMaterial,
  Color3,
  Vector3,
} from '@babylonjs/core'

const BULLET_SPEED = 25
const BULLET_LIFETIME = 3 // seconds

export class BulletEntity {
  private mesh: Mesh
  private velocity: Vector3
  private lifetime = 0
  public isDead = false

  constructor(scene: Scene, origin: Vector3, target: Vector3) {
    this.mesh = MeshBuilder.CreateSphere(
      `bullet-${Date.now()}-${Math.random()}`,
      { diameter: 0.15 },
      scene,
    )
    this.mesh.position = origin.clone()

    const mat = new StandardMaterial(`bullet-mat-${Date.now()}`, scene)
    mat.diffuseColor = new Color3(1, 0.9, 0.2)
    mat.emissiveColor = new Color3(1, 0.7, 0)
    mat.disableLighting = true
    this.mesh.material = mat

    const dir = target.subtract(origin).normalize()
    this.velocity = dir.scale(BULLET_SPEED)
  }

  update(deltaTime: number): boolean {
    if (this.isDead) return true

    this.lifetime += deltaTime
    if (this.lifetime >= BULLET_LIFETIME) {
      this.destroy()
      return true
    }

    this.mesh.position.addInPlace(this.velocity.scale(deltaTime))
    return false
  }

  getPosition(): Vector3 {
    return this.mesh.position
  }

  getMesh(): Mesh {
    return this.mesh
  }

  destroy() {
    if (!this.isDead) {
      this.isDead = true
      this.mesh.dispose()
    }
  }
}
