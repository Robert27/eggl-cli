import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'eggl',
  description: 'A general-purpose helper CLI for dev workflow, cloud, and everyday tasks.',
  lang: 'en-US',
  base: '/',
  cleanUrls: true,
  lastUpdated: true,
  appearance: 'dark',
  head: [
    ['link', { rel: 'icon', href: '/favicon.svg', type: 'image/svg+xml' }],
    ['meta', { name: 'theme-color', content: '#F97316' }],
    ['meta', { property: 'og:title', content: 'eggl — developer helper CLI' }],
    ['meta', { property: 'og:description', content: 'File hygiene, env switching, port-forwards, and more — one CLI for your daily workflow.' }],
    ['link', { rel: 'preconnect', href: 'https://fonts.googleapis.com' }],
    ['link', { rel: 'preconnect', href: 'https://fonts.gstatic.com', crossorigin: '' }],
    ['link', { rel: 'stylesheet', href: 'https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;600&family=Outfit:wght@300;400;500;600;700;800&display=swap' }],
  ],
  themeConfig: {
    logo: '/logo.svg',
    siteTitle: 'eggl',
    nav: [
      { text: 'Guide', link: '/guide/getting-started', activeMatch: '/guide/' },
      { text: 'Commands', link: '/commands/', activeMatch: '/commands/' },
      {
        text: 'GitHub',
        link: 'https://github.com/Robert27/eggl-cli',
      },
    ],
    sidebar: {
      '/guide/': [
        {
          text: 'Introduction',
          items: [
            { text: 'Getting Started', link: '/guide/getting-started' },
            { text: 'Installation', link: '/guide/installation' },
            { text: 'Configuration', link: '/guide/configuration' },
            { text: 'Shell Completion', link: '/guide/completion' },
          ],
        },
      ],
      '/commands/': [
        {
          text: 'Overview',
          items: [{ text: 'All Commands', link: '/commands/' }],
        },
        {
          text: 'File & Git',
          items: [
            { text: 'dedash', link: '/commands/dedash' },
            { text: 'eol', link: '/commands/eol' },
            { text: 'empty', link: '/commands/empty' },
          ],
        },
        {
          text: 'Cloud & Dev',
          items: [
            { text: 'env', link: '/commands/env' },
            { text: 'pf', link: '/commands/pf' },
            { text: 'kill', link: '/commands/kill' },
          ],
        },
        {
          text: 'Utilities',
          items: [
            { text: 'doctor', link: '/commands/doctor' },
            { text: 'version', link: '/commands/version' },
            { text: 'completion', link: '/commands/completion' },
          ],
        },
      ],
    },
    socialLinks: [
      { icon: 'github', link: 'https://github.com/Robert27/eggl-cli' },
    ],
    search: {
      provider: 'local',
    },
    footer: {
      message:
        'Released under the <a href="https://github.com/Robert27/eggl-cli/blob/main/LICENSE">MIT License</a>.',
      copyright:
        'Copyright © <a href="https://github.com/Robert27">Robert Eggl</a> · <a href="/imprint">Imprint</a>',
    },
  },
})
