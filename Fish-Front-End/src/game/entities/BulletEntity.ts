const BULLET_SPEED = 620 // px/s
const BULLET_MAX_LIFE = 1.8 // seconds

export class BulletEntity {
  public x: number
  public y: number
  public isDead = false

  private vx: number
  private vy: number
  private lifetime = 0
  private trail: Array<{ x: number; y: number }> = []

  constructor(startX: number, startY: number, targetX: number, targetY: number) {
    const dx = targetX - startX
    const dy = targetY - startY
    const dist = Math.sqrt(dx * dx + dy * dy) || 1
    this.x = startX
    this.y = startY
    this.vx = (dx / dist) * BULLET_SPEED
    this.vy = (dy / dist) * BULLET_SPEED
  }

  update(dt: number) {
    if (this.isDead) return

    this.trail.push({ x: this.x, y: this.y })
    if (this.trail.length > 7) this.trail.shift()

    this.x += this.vx * dt
    this.y += this.vy * dt
    this.lifetime += dt

    if (this.lifetime >= BULLET_MAX_LIFE) this.isDead = true
  }

  draw(ctx: CanvasRenderingContext2D) {
    if (this.isDead) return

    // Trail
    for (let i = 0; i < this.trail.length; i++) {
      const t = i / this.trail.length
      ctx.beginPath()
      ctx.arc(this.trail[i].x, this.trail[i].y, 3 * t, 0, Math.PI * 2)
      ctx.fillStyle = `rgba(251,191,36,${t * 0.35})`
      ctx.fill()
    }

    // Glow halo
    const grd = ctx.createRadialGradient(this.x, this.y, 0, this.x, this.y, 12)
    grd.addColorStop(0, 'rgba(255,220,50,0.7)')
    grd.addColorStop(1, 'rgba(255,140,0,0)')
    ctx.beginPath()
    ctx.arc(this.x, this.y, 12, 0, Math.PI * 2)
    ctx.fillStyle = grd
    ctx.fill()

    // Core
    ctx.beginPath()
    ctx.arc(this.x, this.y, 4, 0, Math.PI * 2)
    ctx.fillStyle = '#fef08a'
    ctx.fill()
  }

  destroy() {
    this.isDead = true
  }
}
