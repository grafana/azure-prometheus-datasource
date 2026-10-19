import CopyWebpackPlugin from 'copy-webpack-plugin';
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
    plugins: [
      new CopyWebpackPlugin({
        patterns: [
          { from: '../pkg/schema/dsconfig.json', to: './schema/dsconfig.json' },
          { from: '../pkg/schema/schema.gen.json', to: './schema/v0alpha1.json' },
          { from: '../pkg/schema/settings.gen.json', to: './schema/v0alpha1/settings.json' },
          { from: '../pkg/schema/settings.examples.gen.json', to: './schema/v0alpha1/settings.examples.json' },
        ],
      }),
    ],
  });
};

export default config;
