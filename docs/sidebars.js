module.exports = {
  docs: [
    {
      type: "category",
      label: "💡 Getting Started",
      collapsed: false,
      items: [
        "introduction",
        "upgrading",
        {
          type: "category",
          label: "🚀 Installation",
          collapsed: false,
          items: ["windows", "macos", "linux"],
        },
      ],
    },
    {
      type: "category",
      label: "⚙️ Configuration",
      items: [
        "config-overview",
        "config-block",
        "config-segment",
        "config-sample",
        "config-title",
        "config-colors",
        "config-templates",
        "config-transient",
        "config-tooltips",
        "config-fonts"
      ],
    },
    {
      type: "category",
      label: "🌟 Segments",
      collapsed: true,
      items: [
        "angular",
        "aws",
        "az",
        "azfunc",
        "battery",
        "brewfather",
        "command",
        "crystal",
        "dart",
        "dotnet",
        "executiontime",
        "exit",
        "git",
        "poshgit",
        "golang",
        "ipify",
        "java",
        "julia",
        "kubectl",
        "nbgv",
        "nightscout",
        "node",
        "os",
        "owm",
        "path",
        "php",
        "plastic",
        "python",
        "root",
        "ruby",
        "rust",
        "session",
        "shell",
        "spotify",
        "strava",
        "sysinfo",
        "terraform",
        "text",
        "time",
        "wakatime",
        "wifi",
        "winreg",
        "ytm",
      ],
    },
    {
      type: "category",
      label: "🙋🏾‍♀️ Contributing",
      collapsed: true,
      items: [
        "contributing_started",
        "contributing_segment",
        "contributing_git",
        "contributing_plastic",
      ],
    },
    "themes",
    "share",
    "faq",
    "contributors",
  ],
};
