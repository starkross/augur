import {themes as prismThemes} from 'prism-react-renderer';
import type {Config} from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

const config: Config = {
  title: 'augur',
  tagline: 'A fast, opinionated linter for OpenTelemetry Collector configurations',
  favicon: 'img/favicon.ico',

  future: {
    v4: true,
  },

  url: 'https://starkross.github.io',
  baseUrl: '/augur/',

  organizationName: 'starkross',
  projectName: 'augur',

  onBrokenLinks: 'throw',

  markdown: {
    hooks: {
      onBrokenMarkdownLinks: 'warn',
    },
  },

  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      {
        docs: {
          sidebarPath: './sidebars.ts',
          routeBasePath: 'docs',
          editUrl: 'https://github.com/starkross/augur/tree/main/docs/site/',
        },
        blog: false,
        theme: {
          customCss: './src/css/custom.css',
        },
      } satisfies Preset.Options,
    ],
  ],

  themeConfig: {
    image: 'img/augur-social-card.jpg',
    colorMode: {
      respectPrefersColorScheme: true,
    },
    navbar: {
      title: 'augur',
      logo: {
        alt: 'augur logo',
        src: 'img/logo.png',
      },
      items: [
        {
          type: 'docSidebar',
          sidebarId: 'docsSidebar',
          position: 'left',
          label: 'Docs',
        },
        {
          to: '/docs/rules',
          label: 'Rules',
          position: 'left',
        },
        {
          href: 'https://github.com/starkross/augur',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Docs',
          items: [
            {label: 'Introduction', to: '/docs/intro'},
            {label: 'Install', to: '/docs/install'},
            {label: 'Quick start', to: '/docs/quick-start'},
            {label: 'Rules', to: '/docs/rules'},
          ],
        },
        {
          title: 'Project',
          items: [
            {label: 'GitHub', href: 'https://github.com/starkross/augur'},
            {label: 'Issues', href: 'https://github.com/starkross/augur/issues'},
            {label: 'Releases', href: 'https://github.com/starkross/augur/releases'},
            {label: 'pkg.go.dev', href: 'https://pkg.go.dev/github.com/starkross/augur'},
          ],
        },
        {
          title: 'More',
          items: [
            {label: 'OpenTelemetry Collector', href: 'https://opentelemetry.io/docs/collector/'},
            {label: 'OPA / Rego', href: 'https://www.openpolicyagent.org/'},
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} augur contributors. Built with Docusaurus.`,
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
      additionalLanguages: ['bash', 'yaml', 'json', 'go'],
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
