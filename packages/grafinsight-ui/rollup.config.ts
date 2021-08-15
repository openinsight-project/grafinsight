import resolve from '@rollup/plugin-node-resolve';
import commonjs from '@rollup/plugin-commonjs';
import image from '@rollup/plugin-image';
import { terser } from 'rollup-plugin-terser';

const pkg = require('./package.json');

const libraryName = pkg.name;

const buildCjsPackage = ({ env }) => {
  return {
    input: `compiled/index.js`,
    output: [
      {
        dir: 'dist',
        name: libraryName,
        format: 'cjs',
        sourcemap: true,
        strict: false,
        exports: 'named',
        chunkFileNames: `[name].${env}.js`,
        globals: {
          react: 'React',
          'prop-types': 'PropTypes',
        },
      },
    ],
    external: [
      'react',
      'react-dom',
      '@grafinsight/data',
      '@grafinsight/e2e-selectors',
      'moment',
      'monaco-editor', // Monaco should not be used directly
      'monaco-editor/esm/vs/editor/editor.api', // Monaco should not be used directly
      'react-monaco-editor',
      'jquery', // required to use jquery.plot, which is assigned externally
    ],
    plugins: [
      commonjs({
        include: /node_modules/,
      }),
      resolve(),
      image(),
      env === 'production' && terser(),
    ],
  };
};
export default [buildCjsPackage({ env: 'development' }), buildCjsPackage({ env: 'production' })];
