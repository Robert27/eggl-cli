<template>
  <footer class="site-footer">
    <div class="site-footer-inner">
      <p class="footer-message">
        Released under the
        <a href="https://github.com/Robert27/eggl-cli/blob/main/LICENSE">MIT License</a>.
      </p>
      <p class="footer-meta">
        <span>Copyright © <a href="https://github.com/Robert27">Robert Eggl</a></span>
        <span class="sep" aria-hidden="true">·</span>
        <a href="/imprint">Imprint</a>
        <span class="sep" aria-hidden="true">·</span>
        <a href="https://github.com/Robert27/eggl-cli">GitHub</a>
      </p>
    </div>
  </footer>
</template>

<style scoped>
.site-footer {
  border-top: 1px solid var(--vp-c-divider);
  background: var(--vp-c-bg);
  padding: 2rem 1.5rem;
}

.site-footer-inner {
  max-width: 1200px;
  margin: 0 auto;
  text-align: center;
}

.footer-message,
.footer-meta {
  margin: 0;
  font-size: 0.875rem;
  line-height: 1.8;
  color: var(--vp-c-text-2);
}

.footer-meta {
  margin-top: 0.35rem;
}

.site-footer a {
  color: var(--vp-c-text-1);
  text-decoration: none;
  transition: color 0.2s;
}

.site-footer a:hover {
  color: var(--vp-c-brand-1);
}

.sep {
  margin: 0 0.5rem;
  color: var(--vp-c-text-3);
}
</style>
