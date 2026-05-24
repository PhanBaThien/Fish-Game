import type { Fish } from '../../types'

const PALETTE = [
  { body: '#22d3ee', shade: '#0891b2' },
  { body: '#fb923c', shade: '#c2410c' },
  { body: '#a78bfa', shade: '#6d28d9' },
  { body: '#4ade80', shade: '#15803d' },
  { body: '#f472b6', shade: '#be185d' },
  { body: '#fbbf24', shade: '#b45309' },
  { body: '#f87171', shade: '#b91c1c' },
]

export class FishEntity {
  public x: number
  public y: number
  public isDead = false
  public currentHealth: number
  public readonly fishData: Fish

  private direction: 1 | -1
  private baseY: number
  private amplitude: number
  private frequency: number
  private speed: number
  private time: number
  private color: { body: string; shade: string }
  readonly size: number // hit-radius

  public onDeath?: (fish: FishEntity) => void

  constructor(fishData: Fish, canvasW: number, canvasH: number, index: number) {
    this.fishData = fishData
    this.currentHealth = fishData.health

    this.direction = index % 2 === 0 ? 1 : -1
    this.baseY = 90 + (index % 8) * ((canvasH - 200) / 8) + (Math.random() - 0.5) * 20
    this.amplitude = 14 + Math.random() * 22
    this.frequency = 0.4 + Math.random() * 0.9
    this.speed = fishData.speed * 55 + 25
    this.time = Math.random() * Math.PI * 2

    this.x = this.direction === 1 ? -80 : canvasW + 80
    this.y = this.baseY

    const ci = Math.min(fishData.reward_multiplier - 1, PALETTE.length - 1)
    this.color = PALETTE[Math.max(0, ci)]

    // size grows with health (log scale), clamped 20-55
    this.size = Math.min(55, Math.max(20, 18 + Math.log(fishData.health + 1) * 9))
  }

  takeDamage(amount: number): boolean {
    if (this.isDead) return false
    this.currentHealth -= amount
    if (this.currentHealth <= 0) {
      this.isDead = true
      this.onDeath?.(this)
      return true
    }
    return false
  }

  update(dt: number, canvasW: number, canvasH: number) {
    if (this.isDead) return
    this.time += dt
    this.x += this.direction * this.speed * dt
    this.y = this.baseY + Math.sin(this.time * this.frequency) * this.amplitude

    if (this.direction === 1 && this.x > canvasW + 100) {
      this.x = -80
      this.baseY = 90 + Math.random() * (canvasH - 200)
      this.time = 0
    } else if (this.direction === -1 && this.x < -100) {
      this.x = canvasW + 80
      this.baseY = 90 + Math.random() * (canvasH - 200)
      this.time = 0
    }
  }

  draw(ctx: CanvasRenderingContext2D) {
    if (this.isDead) return

    ctx.save()
    ctx.translate(this.x, this.y)
    if (this.direction === -1) ctx.scale(-1, 1)

    const s = this.size
    const { body, shade } = this.color

    // Shadow
    ctx.save()
    ctx.globalAlpha = 0.15
    ctx.beginPath()
    ctx.ellipse(s * 0.1, s * 0.75, s * 0.8, s * 0.2, 0, 0, Math.PI * 2)
    ctx.fillStyle = '#000'
    ctx.fill()
    ctx.restore()

    // Tail
    ctx.beginPath()
    ctx.moveTo(-s * 0.85, 0)
    ctx.lineTo(-s * 1.55, -s * 0.55)
    ctx.lineTo(-s * 1.55, s * 0.55)
    ctx.closePath()
    ctx.fillStyle = shade
    ctx.fill()

    // Body
    ctx.beginPath()
    ctx.ellipse(0, 0, s, s * 0.58, 0, 0, Math.PI * 2)
    ctx.fillStyle = body
    ctx.fill()

    // Belly highlight
    ctx.beginPath()
    ctx.ellipse(s * 0.08, s * 0.18, s * 0.48, s * 0.22, 0, 0, Math.PI * 2)
    ctx.fillStyle = 'rgba(255,255,255,0.22)'
    ctx.fill()

    // Dorsal fin
    ctx.beginPath()
    ctx.moveTo(-s * 0.05, -s * 0.54)
    ctx.lineTo(s * 0.25, -s * 0.96)
    ctx.lineTo(s * 0.58, -s * 0.54)
    ctx.closePath()
    ctx.fillStyle = shade
    ctx.fill()

    // Pectoral fin
    ctx.beginPath()
    ctx.ellipse(s * 0.15, s * 0.3, s * 0.28, s * 0.13, -0.4, 0, Math.PI * 2)
    ctx.fillStyle = `${shade}cc`
    ctx.fill()

    // Eye
    ctx.beginPath()
    ctx.arc(s * 0.52, -s * 0.1, s * 0.16, 0, Math.PI * 2)
    ctx.fillStyle = '#0f172a'
    ctx.fill()
    ctx.beginPath()
    ctx.arc(s * 0.55, -s * 0.13, s * 0.07, 0, Math.PI * 2)
    ctx.fillStyle = '#fff'
    ctx.fill()

    ctx.restore()

    // Health bar (un-flipped world space)
    this.drawHealthBar(ctx)
  }

  private drawHealthBar(ctx: CanvasRenderingContext2D) {
    const ratio = Math.max(0, this.currentHealth / this.fishData.health)
    const barW = this.size * 2.2
    const barH = 6
    const bx = this.x - barW / 2
    const by = this.y - this.size - 16

    // Shadow
    ctx.fillStyle = 'rgba(0,0,0,0.45)'
    ctx.fillRect(bx - 1, by - 1, barW + 2, barH + 2)

    // Track
    ctx.fillStyle = '#450a0a'
    ctx.fillRect(bx, by, barW, barH)

    // Fill — green → yellow → red
    const r = Math.round((1 - ratio) * 255)
    const g = Math.round(ratio * 210)
    ctx.fillStyle = `rgb(${r},${g},20)`
    ctx.fillRect(bx, by, barW * ratio, barH)
  }
}
