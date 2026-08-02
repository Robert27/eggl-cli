import * as path from "node:path";
import { defineConfig, type UserConfig } from "@rspress/core";
import mermaid from "rspress-plugin-mermaid";

const config: UserConfig = {
  root: path.join(__dirname, "docs"),
  title: "eggl-cli",
  description:
    "eggl-cli is a small command-line toolkit for development, Kubernetes, and everyday workflows.",
  icon: "/favicon.svg",
  logoText: "eggl-cli",
  base: "/",
  lang: "en",
  head: [
    ["meta", { name: "author", content: "Robert Eggl" }],
    ["link", { rel: "author", href: "https://eggl.dev" }],
  ],
  plugins: [
    mermaid({
      mermaidConfig: {
        theme: "base",
        themeVariables: {
          primaryColor: "#fff0e6",
          primaryTextColor: "#17191c",
          primaryBorderColor: "#c2410c",
          lineColor: "#5d646b",
          secondaryColor: "#ebe8e2",
          tertiaryColor: "#f7f5f1",
          fontFamily: "Instrument Sans, ui-sans-serif, system-ui, sans-serif",
        },
        flowchart: { curve: "basis" },
      },
    }),
  ],
  markdown: {
    link: {
      checkDeadLinks: true,
    },
  },
  route: {
    cleanUrls: true,
  },
  globalStyles: path.join(__dirname, "styles/index.css"),
  themeConfig: {
    socialLinks: [
      {
        icon: "github",
        mode: "link",
        content: "https://github.com/roberteggl/eggl-cli",
      },
    ],
    editLink: {
      docRepoBaseUrl:
        "https://github.com/roberteggl/eggl-cli/tree/main/docs",
    },
    lastUpdated: true,
    enableScrollToTop: true,
    footer: {
      message:
        'Built by <a href="https://eggl.dev">Robert Eggl</a>.<br />Released under the <a href="https://github.com/roberteggl/eggl-cli/blob/main/LICENSE">MIT License</a>.',
    },
  },
  builderConfig: {
    html: {
      tags: [
        {
          tag: "meta",
          attrs: {
            name: "theme-color",
            content: "#c2410c",
          },
        },
      ],
    },
  },
};

export default defineConfig(config);
