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

let _instanceCounter = 0

export class FishEntity {
  public x: number
  public y: number
  public isDead = false
  public readonly instanceId: string
  public readonly killProb: number   // xác suất kill mỗi lần bắn trúng (đã nhân RTP)
  public isFlashing = false
  private flashTimer = 0
  private shakeX = 0                // rung lắc khi bị trúng
  private shakeY = 0
  public readonly fishData: Fish

  private direction: 1 | -1
  private baseY: number
  private amplitude: number
  private frequency: number
  private speed: number
  private time: number
  private color: { body: string; shade: string }
  readonly size: number

  constructor(fishData: Fish, canvasW: number, canvasH: number, index: number, killProb: number) {
    this.fishData = fishData
    this.instanceId = `fish_${++_instanceCounter}`
    this.killProb = killProb

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

  // Bật flash + shake — không set isDead (server mới quyết định)
  takeDamage(_amount: number) {
    if (this.isDead) return
    this.isFlashing = true
    this.flashTimer = 0.20
    // rung ngẫu nhiên nhỏ
    this.shakeX = (Math.random() - 0.5) * 6
    this.shakeY = (Math.random() - 0.5) * 4
  }

  // Gọi khi server xác nhận cá đã chết
  confirmDeath() {
    if (this.isDead) return
    this.isDead = true
  }

  update(dt: number, canvasW: number, canvasH: number) {
    if (this.isDead) return
    this.time += dt
    this.x += this.direction * this.speed * dt
    this.y = this.baseY + Math.sin(this.time * this.frequency) * this.amplitude

    // Flash + shake timer
    if (this.isFlashing) {
      this.flashTimer -= dt
      // shake tắt dần nhanh hơn flash
      this.shakeX *= 0.7
      this.shakeY *= 0.7
      if (this.flashTimer <= 0) {
        this.isFlashing = false
        this.shakeX = 0
        this.shakeY = 0
      }
    }

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
    ctx.translate(this.x + this.shakeX, this.y + this.shakeY)

    // Flash effect: lớp trắng phủ lên cá khi bị trúng đạn
    if (this.isFlashing) {
      const alpha = (this.flashTimer / 0.20) * 0.55
      ctx.save()
      ctx.globalAlpha = alpha
      ctx.beginPath()
      ctx.ellipse(0, 0, this.size * 1.1, this.size * 0.65, 0, 0, Math.PI * 2)
      ctx.fillStyle = '#ffffff'
      ctx.fill()
      ctx.restore()
    }

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

    // Label tỷ lệ (world space, không bị ảnh hưởng bởi flip/translate)
    this.drawStatsLabel(ctx)
  }

  private drawStatsLabel(ctx: CanvasRenderingContext2D) {
    const cx = this.x + this.shakeX
    const cy = this.y + this.shakeY - this.size - 10

    const mult    = this.fishData.reward_multiplier
    const pct     = (this.killProb * 100)
    const probStr = pct < 0.1
      ? pct.toFixed(3) + '%'
      : pct < 1
      ? pct.toFixed(2) + '%'
      : pct.toFixed(1) + '%'

    // Màu xác suất: xanh → vàng → cam → đỏ tuỳ theo độ khó
    const probColor = pct > 10 ? '#4ade80'
      : pct > 1   ? '#facc15'
      : pct > 0.1 ? '#fb923c'
      :              '#f87171'

    ctx.save()
    ctx.font = 'bold 11px system-ui, sans-serif'
    ctx.textAlign = 'center'
    ctx.textBaseline = 'middle'

    const multText = `×${mult}`
    const probText = probStr

    const multW  = ctx.measureText(multText).width + 10
    const probW  = ctx.measureText(probText).width + 10
    const gap    = 4
    const totalW = multW + gap + probW
    const h      = 16
    const r      = 4

    const startX = cx - totalW / 2

    // Nền multiplier badge
    ctx.fillStyle = 'rgba(0,0,0,0.55)'
    this._roundRect(ctx, startX, cy - h / 2, multW, h, r)
    ctx.fill()
    ctx.fillStyle = '#e2e8f0'
    ctx.fillText(multText, startX + multW / 2, cy)

    // Nền prob badge
    ctx.fillStyle = 'rgba(0,0,0,0.55)'
    this._roundRect(ctx, startX + multW + gap, cy - h / 2, probW, h, r)
    ctx.fill()
    ctx.fillStyle = probColor
    ctx.fillText(probText, startX + multW + gap + probW / 2, cy)

    ctx.restore()
  }

  private _roundRect(ctx: CanvasRenderingContext2D, x: number, y: number, w: number, h: number, r: number) {
    ctx.beginPath()
    ctx.moveTo(x + r, y)
    ctx.lineTo(x + w - r, y)
    ctx.arcTo(x + w, y, x + w, y + r, r)
    ctx.lineTo(x + w, y + h - r)
    ctx.arcTo(x + w, y + h, x + w - r, y + h, r)
    ctx.lineTo(x + r, y + h)
    ctx.arcTo(x, y + h, x, y + h - r, r)
    ctx.lineTo(x, y + r)
    ctx.arcTo(x, y, x + r, y, r)
    ctx.closePath()
  }
}
