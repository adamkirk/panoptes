interface Project {
  enabled: boolean;
  // opts will be merged into the playwright project config
  opts?: {};
}

interface Config {
  verifySSL: boolean;
  enabledBrowsers: string[];
  projects: {
    api: Project;
  };
}

const config: Config = {
  verifySSL: false,
  enabledBrowsers: [
    "chromium",
  ],
  projects: {
    api: {
      enabled: true,
      opts: {
        use: {
          baseURL: "http://heimdallr.test",
          ignoreHTTPSErrors: true,
          extraHTTPHeaders: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
          },
        },
      },
    },
  },
};

export default config