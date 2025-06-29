import type { Configuration } from 'webpack';
import { merge } from 'webpack-merge';

import grafanaConfig from './.config/webpack/webpack.config';

const config = async (env): Promise<Configuration> => {
  const baseConfig = await grafanaConfig(env);

  return merge(baseConfig, {
    externals: ['i18next'],
    module: {
      rules: [
        {
          test: /\.(m|c)?js/,
          resolve: {
            fullySpecified: false,
          },
        },
      ],
    },
  });
};

export default config;
