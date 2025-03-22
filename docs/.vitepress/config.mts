import { defineConfig } from "vitepress";

export default defineConfig({
  title: "Quokka",
  titleTemplate: ":title ~ Quokka",
  description: "One-Command, Self-Hosted Mail Server",
  head: [
    ["link", { rel: "icon", type: "image/svg+xml", href: "/images/logo.svg" }],
    ["link", { rel: "icon", type: "image/png", href: "/images/logo.png" }],
  ],
  themeConfig: {
    logo: { src: "/images/logo.svg" },
    nav: [
      {
        text: "Guide",
        link: "/guide/what-is-quokka",
        activeMatch: "/guide/",
      },
    ],
    sidebar: {
      "/guide/": [
        {
          text: "Introduction",
          items: [
            { text: "What is Quokka?", link: "/guide/what-is-quokka" },
            { text: "Getting Started", link: "/guide/getting-started" },
          ],
        },
      ],
    },
    socialLinks: [
      { icon: "github", link: "https://github.com/quokkamail/quokka" },
    ],
    editLink: {
      pattern: "https://github.com/quokkamail/quokka/edit/main/docs/:path",
      text: "Edit this page on GitHub",
    },
    lastUpdated: {
      text: "Updated at",
      formatOptions: {
        dateStyle: "full",
        timeStyle: "medium",
      },
    },
  },
});
