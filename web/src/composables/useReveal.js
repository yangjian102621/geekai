// * 滚动渐入指令：作用于页面根元素，观察其内部所有 [data-reveal] 子元素。
// * 在 <script setup> 中用法：const vReveal = useReveal(); <div v-reveal>

export function useReveal() {
  return {
    mounted(el) {
      const targets = el.querySelectorAll('[data-reveal]')
      if (targets.length === 0) return

      const observer = new IntersectionObserver(
        (entries) => {
          entries.forEach((entry) => {
            if (entry.isIntersecting) {
              entry.target.classList.add('reveal-visible')
              observer.unobserve(entry.target)
            }
          })
        },
        { threshold: 0.1, rootMargin: '0px 0px -60px 0px' },
      )

      targets.forEach((t) => observer.observe(t))
      el.__revealObserver__ = observer
    },
    unmounted(el) {
      el.__revealObserver__?.disconnect()
      delete el.__revealObserver__
    },
  }
}
