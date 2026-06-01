import { FishEntity } from '../entities/FishEntity'
import { BulletEntity } from '../entities/BulletEntity'
import type { Fish } from '../../types'

export interface GameSceneOptions {
  canvas: HTMLCanvasElement
  fishList: Fish[]
  roomRtp: number                                           // RTP phòng dạng decimal (0.0–1.0)
  onHitFish?: (fishId: number, instanceId: string) => void // đạn chạm cá → gửi lên server
  onScore?: (points: number) => void
  onShot?: (x: number, y: number, angle: number) => boolean // trả false → không bắn đạn
}

interface Particle {
  x: number; y: number
  vx: number; vy: number
  life: number; maxLife: number
  color: string; size: number
}

interface Bubble {
  x: number; y: number; r: number
  vy: number; wobble: number; freq: number; t: number
}

export class GameScene {
  private ctx: CanvasRenderingContext2D
  private canvas: HTMLCanvasElement
  private options: GameSceneOptions
  private fishEntities: FishEntity[] = []
  private bullets: BulletEntity[] = []
  private particles: Particle[] = []
  private bubbles: Bubble[] = []

  private mouseX = 0
  private mouseY = 0
  private cannonAngle = -Math.PI / 2
  private isDisposed = false
  private animFrame = 0
  private lastTime = 0
  private bgTime = 0

  constructor(options: GameSceneOptions) {
    this.options = options
    this.canvas = options.canvas
    const ctx = options.canvas.getContext('2d')
    if (!ctx) throw new Error('Cannot get 2D context')
    this.ctx = ctx

    this.resize()
    this.initBubbles()
    this.spawnFish()
    this.setupEvents()
    this.loop(performance.now())
  }

  // ── Setup ─────────────────────────────────────────────────────────────────

  private resize() {
    const parent = this.canvas.parentElement
    const rect = parent?.getBoundingClientRect()
    this.canvas.width = rect?.width ?? window.innerWidth
    this.canvas.height = rect?.height ?? window.innerHeight
  }

  private initBubbles() {
    const { width: w, height: h } = this.canvas
    for (let i = 0; i < 28; i++) {
      this.bubbles.push({
        x: Math.random() * w,
        y: Math.random() * h,
        r: 2 + Math.random() * 6,
        vy: 18 + Math.random() * 28,
        wobble: Math.random() * Math.PI * 2,
        freq: 0.4 + Math.random() * 1.4,
        t: Math.random() * 100,
      })
    }
  }

  private spawnFish() {
    const { fishList } = this.options
    if (!fishList?.length) return
    const count = Math.max(fishList.length * 2, 8)
    for (let i = 0; i < count; i++) {
      this.addFish(fishList[i % fishList.length], i)
    }
  }

  addFish(data: Fish, index: number) {
    const rtp = this.options.roomRtp > 0 && this.options.roomRtp <= 1 ? this.options.roomRtp : 0.90
    const killProb = data.base_prob * rtp
    const fish = new FishEntity(data, this.canvas.width, this.canvas.height, index, killProb)
    this.fishEntities.push(fish)
  }

  // Gọi khi server xác nhận cá đã chết (hit_result.killed = true)
  confirmFishDeath(instanceId: string) {
    const fish = this.fishEntities.find((f) => f.instanceId === instanceId)
    if (!fish || fish.isDead) return

    this.spawnParticles(fish.x, fish.y, fish.size)
    this.options.onScore?.(10)
    fish.confirmDeath()

    const deadX = fish.x
    const deadY = fish.y
    const fishList = this.options.fishList

    setTimeout(() => {
      if (this.isDisposed) return
      const idx = this.fishEntities.findIndex((f) => f.instanceId === instanceId)
      if (idx !== -1) this.fishEntities.splice(idx, 1)
      if (fishList?.length) {
        this.addFish(fishList[Math.floor(Math.random() * fishList.length)], Math.floor(Math.random() * 100))
      }
    }, 600)

    void deadX; void deadY // suppress unused warning
  }

  private spawnParticles(cx: number, cy: number, radius: number) {
    const colors = ['#fbbf24', '#f59e0b', '#ef4444', '#fb923c', '#fff', '#fde68a']
    const count = 18 + Math.round(radius * 0.6)
    for (let i = 0; i < count; i++) {
      const angle = (Math.PI * 2 * i) / count + Math.random() * 0.6
      const speed = 90 + Math.random() * 210
      this.particles.push({
        x: cx, y: cy,
        vx: Math.cos(angle) * speed,
        vy: Math.sin(angle) * speed,
        life: 0,
        maxLife: 0.35 + Math.random() * 0.45,
        color: colors[Math.floor(Math.random() * colors.length)],
        size: 3 + Math.random() * 5,
      })
    }
  }

  // ── Events ────────────────────────────────────────────────────────────────

  private setupEvents() {
    this.canvas.addEventListener('mousemove', this.onMouseMove)
    this.canvas.addEventListener('click', this.onClick)
    window.addEventListener('resize', this.onResize)
  }

  private onMouseMove = (e: MouseEvent) => {
    const rect = this.canvas.getBoundingClientRect()
    this.mouseX = e.clientX - rect.left
    this.mouseY = e.clientY - rect.top
    const cx = this.canvas.width / 2
    const cy = this.canvas.height - 55
    this.cannonAngle = Math.atan2(this.mouseY - cy, this.mouseX - cx)
  }

  private onClick = (e: MouseEvent) => {
    const rect = this.canvas.getBoundingClientRect()
    const tx = e.clientX - rect.left
    const ty = e.clientY - rect.top
    const cx = this.canvas.width / 2
    const cy = this.canvas.height - 55
    const tipDist = 38
    const sx = cx + Math.cos(this.cannonAngle) * tipDist
    const sy = cy + Math.sin(this.cannonAngle) * tipDist
    // onShot trả false (không đủ tiền) → không hiển thị đạn
    const allowed = this.options.onShot?.(tx, ty, this.cannonAngle) ?? true
    if (allowed) {
      this.bullets.push(new BulletEntity(sx, sy, tx, ty))
    }
  }

  private onResize = () => {
    this.resize()
  }

  // ── Game loop ─────────────────────────────────────────────────────────────

  private loop = (now: number) => {
    if (this.isDisposed) return
    const dt = Math.min((now - this.lastTime) / 1000, 0.05)
    this.lastTime = now
    this.bgTime += dt

    this.update(dt)
    this.draw()
    this.animFrame = requestAnimationFrame(this.loop)
  }

  private update(dt: number) {
    const { width: w, height: h } = this.canvas

    for (const fish of this.fishEntities) fish.update(dt, w, h)

    for (const b of this.bullets) b.update(dt)
    this.bullets = this.bullets.filter((b) => !b.isDead)

    for (const p of this.particles) {
      p.life += dt
      p.x += p.vx * dt
      p.y += p.vy * dt
      p.vy += 180 * dt
    }
    this.particles = this.particles.filter((p) => p.life < p.maxLife)

    for (const b of this.bubbles) {
      b.t += dt
      b.y -= b.vy * dt
      b.x += Math.sin(b.t * b.freq) * 0.6
      if (b.y < -15) { b.y = h + 10; b.x = Math.random() * w }
    }

    // Collision: đạn chạm cá → flash effect + gửi server, không tự quyết định death
    for (const bullet of this.bullets) {
      if (bullet.isDead) continue
      for (const fish of this.fishEntities) {
        if (fish.isDead) continue
        const dx = bullet.x - fish.x
        const dy = bullet.y - fish.y
        if (Math.sqrt(dx * dx + dy * dy) < fish.size * 0.88) {
          fish.takeDamage(10)     // visual flash effect
          bullet.destroy()
          this.options.onHitFish?.(fish.fishData.id, fish.instanceId) // server quyết định kill
          break
        }
      }
    }
  }

  // ── Drawing ───────────────────────────────────────────────────────────────

  private draw() {
    const ctx = this.ctx
    const { width: w, height: h } = this.canvas

    // Background
    const bg = ctx.createLinearGradient(0, 0, 0, h)
    bg.addColorStop(0, '#010d1f')
    bg.addColorStop(0.45, '#041630')
    bg.addColorStop(1, '#062040')
    ctx.fillStyle = bg
    ctx.fillRect(0, 0, w, h)

    this.drawCaustics(ctx, w, h)
    this.drawBubbles(ctx)
    this.drawSeabed(ctx, w, h)

    for (const fish of this.fishEntities) fish.draw(ctx)
    for (const b of this.bullets) b.draw(ctx)

    this.drawParticles(ctx)
    this.drawCannon(ctx, w, h)
    this.drawCrosshair(ctx)
  }

  private drawCaustics(ctx: CanvasRenderingContext2D, w: number, h: number) {
    ctx.save()
    ctx.globalAlpha = 0.035
    for (let i = 0; i < 6; i++) {
      const x = w * (i / 6) + Math.sin(this.bgTime * 0.28 + i * 1.1) * 35
      const grad = ctx.createLinearGradient(x, 0, x + 35, h * 0.65)
      grad.addColorStop(0, 'rgba(100,200,255,1)')
      grad.addColorStop(1, 'rgba(0,50,140,0)')
      ctx.fillStyle = grad
      ctx.beginPath()
      ctx.moveTo(x, 0)
      ctx.lineTo(x + 18 + Math.sin(this.bgTime * 0.45 + i) * 12, h * 0.65)
      ctx.lineTo(x + 36, 0)
      ctx.fill()
    }
    ctx.restore()
  }

  private drawBubbles(ctx: CanvasRenderingContext2D) {
    ctx.save()
    for (const b of this.bubbles) {
      ctx.globalAlpha = 0.22
      ctx.beginPath()
      ctx.arc(b.x, b.y, b.r, 0, Math.PI * 2)
      ctx.strokeStyle = 'rgba(160,230,255,0.9)'
      ctx.lineWidth = 1
      ctx.stroke()
      ctx.globalAlpha = 0.06
      ctx.fillStyle = 'rgba(160,230,255,1)'
      ctx.fill()
    }
    ctx.restore()
  }

  private drawSeabed(ctx: CanvasRenderingContext2D, w: number, h: number) {
    // Sand gradient
    const sand = ctx.createLinearGradient(0, h - 65, 0, h)
    sand.addColorStop(0, 'rgba(8,35,75,0)')
    sand.addColorStop(1, 'rgba(12,45,90,0.85)')
    ctx.fillStyle = sand
    ctx.fillRect(0, h - 65, w, 65)

    // Seaweed clusters
    const positions = [0.06, 0.18, 0.42, 0.68, 0.82, 0.94]
    for (const p of positions) this.drawSeaweed(ctx, w * p, h)
  }

  private drawSeaweed(ctx: CanvasRenderingContext2D, x: number, h: number) {
    ctx.save()
    ctx.strokeStyle = 'rgba(22,101,52,0.55)'
    ctx.lineWidth = 4
    ctx.lineCap = 'round'
    ctx.lineJoin = 'round'
    ctx.beginPath()
    ctx.moveTo(x, h)
    const segs = 5 + Math.floor(Math.random() * 2)
    for (let i = 1; i <= segs; i++) {
      const sway = Math.sin(this.bgTime * 0.75 + x * 0.02 + i * 0.9) * 7 * (i / segs)
      ctx.lineTo(x + sway, h - i * 20)
    }
    ctx.stroke()
    ctx.restore()
  }

  private drawParticles(ctx: CanvasRenderingContext2D) {
    for (const p of this.particles) {
      const alpha = 1 - p.life / p.maxLife
      ctx.save()
      ctx.globalAlpha = alpha
      ctx.beginPath()
      ctx.arc(p.x, p.y, p.size * alpha, 0, Math.PI * 2)
      ctx.fillStyle = p.color
      ctx.fill()
      ctx.restore()
    }
  }

  private drawCannon(ctx: CanvasRenderingContext2D, w: number, h: number) {
    const cx = w / 2
    const cy = h - 55
    const barrelLen = 38
    const barrelW = 13

    ctx.save()
    ctx.translate(cx, cy)

    // Outer glow ring
    ctx.beginPath()
    ctx.arc(0, 0, 33, 0, Math.PI * 2)
    ctx.strokeStyle = `rgba(14,165,233,${0.2 + Math.sin(this.bgTime * 2) * 0.08})`
    ctx.lineWidth = 5
    ctx.stroke()

    // Base platform
    ctx.beginPath()
    ctx.arc(0, 0, 27, 0, Math.PI * 2)
    const baseGrd = ctx.createRadialGradient(-4, -4, 2, 0, 0, 27)
    baseGrd.addColorStop(0, '#475569')
    baseGrd.addColorStop(1, '#1e293b')
    ctx.fillStyle = baseGrd
    ctx.fill()
    ctx.strokeStyle = '#38bdf8'
    ctx.lineWidth = 2
    ctx.stroke()

    // Barrel
    ctx.rotate(this.cannonAngle)
    const barrelGrd = ctx.createLinearGradient(0, -barrelW / 2, 0, barrelW / 2)
    barrelGrd.addColorStop(0, '#64748b')
    barrelGrd.addColorStop(0.5, '#94a3b8')
    barrelGrd.addColorStop(1, '#334155')
    ctx.fillStyle = barrelGrd
    ctx.strokeStyle = '#38bdf8'
    ctx.lineWidth = 1.5

    // Barrel body
    ctx.beginPath()
    ctx.rect(2, -barrelW / 2, barrelLen, barrelW)
    ctx.fill()
    ctx.stroke()

    // Tip ring
    ctx.fillStyle = '#38bdf8'
    ctx.fillRect(barrelLen - 2, -barrelW / 2 - 2, 8, barrelW + 4)

    // Center knob
    ctx.restore()
    ctx.save()
    ctx.translate(cx, cy)
    ctx.beginPath()
    ctx.arc(0, 0, 8, 0, Math.PI * 2)
    ctx.fillStyle = '#38bdf8'
    ctx.fill()

    ctx.restore()
  }

  private drawCrosshair(ctx: CanvasRenderingContext2D) {
    const { x, y } = { x: this.mouseX, y: this.mouseY }
    const arm = 11

    ctx.save()
    ctx.strokeStyle = 'rgba(255,255,255,0.65)'
    ctx.lineWidth = 1.5
    ctx.setLineDash([3, 3])
    ctx.beginPath()
    ctx.moveTo(x - arm, y); ctx.lineTo(x + arm, y)
    ctx.moveTo(x, y - arm); ctx.lineTo(x, y + arm)
    ctx.stroke()

    ctx.setLineDash([])
    ctx.beginPath()
    ctx.arc(x, y, 6, 0, Math.PI * 2)
    ctx.strokeStyle = 'rgba(251,191,36,0.85)'
    ctx.lineWidth = 1.5
    ctx.stroke()
    ctx.restore()
  }

  // ── Cleanup ───────────────────────────────────────────────────────────────

  dispose() {
    this.isDisposed = true
    cancelAnimationFrame(this.animFrame)
    this.canvas.removeEventListener('mousemove', this.onMouseMove)
    this.canvas.removeEventListener('click', this.onClick)
    window.removeEventListener('resize', this.onResize)
  }
}
