import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'viewmd',
  description: 'A minimal file viewer for the browser',
  base: '/viewmd/',

  head: [
    ['meta', { name: 'theme-color', content: '#4a9e6e' }],
  ],

  themeConfig: {
    nav: [
      { text: 'Guide', link: '/getting-started/' },
      { text: 'Reference', link: '/reference/cli' },
    ],

    sidebar: [
      {
        text: 'Getting Started',
        items: [
          { text: 'Introduction', link: '/getting-started/' },
          { text: 'Installation', link: '/getting-started/installation' },
          { text: 'Quick Start', link: '/getting-started/quick-start' },
        ],
      },
      {
        text: 'Guide',
        items: [
          { text: 'Viewing Files', link: '/guide/viewing-files' },
          { text: 'Directory Browsing', link: '/guide/directory-browsing' },
          { text: 'Search', link: '/guide/search' },
          { text: 'Keyboard Shortcuts', link: '/guide/keyboard-shortcuts' },
        ],
      },
      {
        text: 'Reference',
        items: [
          { text: 'CLI Flags', link: '/reference/cli' },
          { text: 'Configuration', link: '/reference/configuration' },
          { text: 'Supported File Types', link: '/reference/file-types' },
        ],
      },
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/driangle/viewmd' },
    ],

    editLink: {
      pattern: 'https://github.com/driangle/viewmd/edit/main/docs/:path',
      text: 'Edit this page on GitHub',
    },

    search: {
      provider: 'local',
    },

    footer: {
      message: 'Released under the MIT License.',
    },
  },
})
