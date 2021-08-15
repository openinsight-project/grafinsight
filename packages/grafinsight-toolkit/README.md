> **WARNING: @grafinsight/toolkit is currently in BETA**.

# grafinsight-toolkit

grafinsight-toolkit is a CLI that enables efficient development of GrafInsight plugins. We want to help our community focus on the core value of their plugins rather than all the setup required to develop them.

## Getting started

Set up a new plugin with `grafinsight-toolkit plugin:create` command:

```sh
npx @grafinsight/toolkit plugin:create my-grafinsight-plugin
cd my-grafinsight-plugin
yarn install
yarn dev
```

## Update your plugin to use grafinsight-toolkit

Follow the steps below to start using grafinsight-toolkit in your existing plugin.

1. Add `@grafinsight/toolkit` package to your project by running `yarn add @grafinsight/toolkit` or `npm install @grafinsight/toolkit`.
2. Create `tsconfig.json` file in the root dir of your plugin and paste the code below:

```json
{
  "extends": "./node_modules/@grafinsight/toolkit/src/config/tsconfig.plugin.json",
  "include": ["src", "types"],
  "compilerOptions": {
    "rootDir": "./src",
    "baseUrl": "./src",
    "typeRoots": ["./node_modules/@types"]
  }
}
```

3. Create `.prettierrc.js` file in the root dir of your plugin and paste the code below:

```js
module.exports = {
  ...require('./node_modules/@grafinsight/toolkit/src/config/prettier.plugin.config.json'),
};
```

4. In your `package.json` file add following scripts:

```json
"scripts": {
  "build": "grafinsight-toolkit plugin:build",
  "test": "grafinsight-toolkit plugin:test",
  "dev": "grafinsight-toolkit plugin:dev",
  "watch": "grafinsight-toolkit plugin:dev --watch"
},
```

## Usage

With grafinsight-toolkit, we give you a CLI that addresses common tasks performed when working on GrafInsight plugin:

- `grafinsight-toolkit plugin:create`
- `grafinsight-toolkit plugin:dev`
- `grafinsight-toolkit plugin:test`
- `grafinsight-toolkit plugin:build`
- `grafinsight-toolkit plugin:sign`

### Create your plugin

`grafinsight-toolkit plugin:create plugin-name`

This command creates a new GrafInsight plugin from template.

If `plugin-name` is provided, then the template is downloaded to `./plugin-name` directory. Otherwise, it will be downloaded to the current directory.

### Develop your plugin

`grafinsight-toolkit plugin:dev`

This command creates a development build that's easy to play with and debug using common browser tooling.

Available options:

- `-w`, `--watch` - run development task in a watch mode

### Test your plugin

`grafinsight-toolkit plugin:test`

This command runs Jest against your codebase.

Available options:

- `--watch` - Runs tests in interactive watch mode.
- `--coverage` - Reports code coverage.
- `-u`, `--updateSnapshot` - Performs snapshots update.
- `--testNamePattern=<regex>` - Runs test with names that match provided regex (https://jestjs.io/docs/en/cli#testnamepattern-regex).
- `--testPathPattern=<regex>` - Runs test with paths that match provided regex (https://jestjs.io/docs/en/cli#testpathpattern-regex).
- `--maxWorkers=<num>|<string>` - Limit number of Jest workers spawned (https://jestjs.io/docs/en/cli#--maxworkersnumstring)

### Build your plugin

`grafinsight-toolkit plugin:build`

This command creates a production-ready build of your plugin.

Available options:

- `--coverage` - Reports code coverage after the test step of the build.
- `--preserveConsole` - Preserves console statements in the code.

### Sign your plugin

`grafinsight-toolkit plugin:sign`

This command creates a signed MANIFEST.txt file which GrafInsight uses to validate the integrity of the plugin.

Available options:

- `--signatureType` - The [type of Signature](https://grafinsight.com/legal/plugins/) you are generating: `private`, `community` or `commercial`
- `--rootUrls` - For private signatures, a list of the GrafInsight instance URLs that the plugin will be used on

To generate a signature, you will need to sign up for a free account on https://grafinsight.com, create an API key with the Plugin Publisher role, and pass that in the `GrafInsight_API_KEY` environment variable.

## FAQ

### Which version of grafinsight-toolkit should I use?

See [GrafInsight packages versioning guide](https://github.com/grafinsight/grafinsight/blob/master/packages/README.md#versioning).

### What tools does grafinsight-toolkit use?

grafinsight-toolkit comes with TypeScript, ESLint, Prettier, Jest, CSS and SASS support.

### How to start using grafinsight-toolkit in my plugin?

See [Updating your plugin to use grafinsight-toolkit](#updating-your-plugin-to-use-grafinsight-toolkit).

### Can I use TypeScript to develop GrafInsight plugins?

Yes! grafinsight-toolkit supports TypeScript by default.

### How can I test my plugin?

grafinsight-toolkit comes with Jest as a test runner.

Internally at GrafInsight we use Enzyme. If you are developing React plugin and you want to configure Enzyme as a testing utility, then you need to configure `enzyme-adapter-react`. To do so, create `<YOUR_PLUGIN_DIR>/config/jest-setup.ts` file that will provide necessary setup. Copy the following code into that file to get Enzyme working with React:

```ts
import { configure } from 'enzyme';
import Adapter from 'enzyme-adapter-react-16';

configure({ adapter: new Adapter() });
```

You can also set up Jest with shims of your needs by creating `jest-shim.ts` file in the same directory: `<YOUR_PLUGIN_DIR_>/config/jest-shim.ts`

### Can I provide custom setup for Jest?

You can provide custom Jest configuration with a `package.json` file. For more details, see [Jest docs](https://jest-bot.github.io/jest/docs/configuration.html).

Currently we support following Jest configuration properties:

- [`snapshotSerializers`](https://jest-bot.github.io/jest/docs/configuration.html#snapshotserializers-array-string)
- [`moduleNameMapper`](https://jestjs.io/docs/en/configuration#modulenamemapper-object-string-string)

### How can I customize Webpack rules or plugins?

You can provide your own `webpack.config.js` file that exports a `getWebpackConfig` function. We recommend that you extend the standard configuration, but you are free to create your own:

```js
const CustomPlugin = require('custom-plugin');

module.exports.getWebpackConfig = (config, options) => ({
  ...config,
  plugins: [...config.plugins, new CustomPlugin()],
});
```

### How can I style my plugin?

We support pure CSS, SASS, and CSS-in-JS approach (via [Emotion](https://emotion.sh/)).

#### Single CSS or SASS file

Create your CSS or SASS file and import it in your plugin entry point (typically `module.ts`):

```ts
import 'path/to/your/css_or_sass';
```

The styles will be injected via `style` tag during runtime.

> **Note:** that imported static assets will be inlined as base64 URIs. _This can be subject of change in the future!_

#### Theme-specific stylesheets

If you want to provide different stylesheets for dark/light theme, then create `dark.[css|scss]` and `light.[css|scss]` files in the `src/styles` directory of your plugin. grafinsight-toolkit generates theme-specific stylesheets that are stored in `dist/styles` directory.

In order for GrafInsight to pick up your theme stylesheets, you need to use `loadPluginCss` from `@grafinsight/runtime` package. Typically you would do that in the entry point of your plugin:

```ts
import { loadPluginCss } from '@grafinsight/runtime';

loadPluginCss({
  dark: 'plugins/<YOUR-PLUGIN-ID>/styles/dark.css',
  light: 'plugins/<YOUR-PLUGIN-ID>/styles/light.css',
});
```

You must add `@grafinsight/runtime` to your plugin dependencies by running `yarn add @grafinsight/runtime` or `npm install @grafinsight/runtime`.

> **Note:** that in this case static files (png, svg, json, html) are all copied to dist directory when the plugin is bundled. Relative paths to those files does not change!

#### Emotion

Starting from GrafInsight 6.2 _our suggested way_ for styling plugins is by using [Emotion](https://emotion.sh). It's a CSS-in-JS library that we use internally at GrafInsight. The biggest advantage of using Emotion is that you can access GrafInsight Theme variables.

To start using Emotion, you first must add it to your plugin dependencies:

```
  yarn add "emotion"@10.0.27
```

Then, import `css` function from Emotion:

```ts
import { css } from 'emotion';
```

Now you are ready to implement your styles:

```tsx
const MyComponent = () => {
  return (
    <div
      className={css`
        background: red;
      `}
    />
  );
};
```

To learn more about using GrafInsight theme please refer to [Theme usage guide](https://github.com/grafinsight/grafinsight/blob/master/style_guides/themes.md#react)

> We do not support Emotion's `css` prop. Use className instead!

### Can I adjust TypeScript configuration to suit my needs?

Yes! However, it's important that your `tsconfig.json` file contains the following lines:

```json
{
  "extends": "./node_modules/@grafinsight/toolkit/src/config/tsconfig.plugin.json",
  "include": ["src"],
  "compilerOptions": {
    "rootDir": "./src",
    "typeRoots": ["./node_modules/@types"]
  }
}
```

### Can I adjust ESLint configuration to suit my needs?

grafinsight-toolkit comes with [default config for ESLint](https://github.com/grafinsight/grafinsight/blob/master/packages/grafinsight-toolkit/src/config/eslint.plugin.json). For now, there is no way to customise ESLint config.

### How is Prettier integrated into grafinsight-toolkit workflow?

When building plugin with [`grafinsight-toolkit plugin:build`](#building-plugin) task, grafinsight-toolkit performs Prettier check. If the check detects any Prettier issues, the build will not pass. To avoid such situation we suggest developing plugin with [`grafinsight-toolkit plugin:dev --watch`](#developing-plugin) task running. This task tries to fix Prettier issues automatically.

### My editor does not respect Prettier config, what should I do?

In order for your editor to pick up our Prettier config you need to create `.prettierrc.js` file in the root directory of your plugin with following content:

```js
module.exports = {
  ...require('./node_modules/@grafinsight/toolkit/src/config/prettier.plugin.config.json'),
};
```

### How do I add third-party dependencies that are not npm packages?

Put them in the `static` directory in the root of your project. The `static` directory is copied when the plugin is built.

### I am getting this message when I run yarn install: `Request failed \"404 Not Found\"`

If you are using version `canary`, this error occurs because a `canary` release unpublishes previous versions leaving `yarn.lock` outdated. Remove `yarn.lock` and run `yarn install` again.

### I am getting this message when I run my plugin: `Unable to dynamically transpile ES module A loader plugin needs to be configured via SystemJS.config({ transpiler: 'transpiler-module' }).`

This error occurs when you bundle your plugin using the `grafinsight-toolkit plugin:dev` task and your code comments include ES2016 code.

There are two issues at play:

- The `grafinsight-toolkit plugin:dev` task does not remove comments from your bundled package.
- GrafInsight does not support [ES modules](https://hacks.mozilla.org/2018/03/es-modules-a-cartoon-deep-dive/).

If your comments include ES2016 code, then SystemJS v0.20.19, which GrafInsight uses internally to load plugins, interprets your code as an ESM and fails.

To fix this error, remove the ES2016 code from your comments.

## Contribute to grafinsight-toolkit

You can contribute to grafinsight-toolkit by helping develop it or by debugging it.

### Develop grafinsight-toolkit

Typically plugins should be developed using the `@grafinsight/toolkit` installed from npm. However, when working on the toolkit, you might want to use the local version. Follow the steps below to develop with a local version:

1. Clone [GrafInsight repository](https://github.com/grafinsight/grafinsight).
2. Navigate to the directory you have cloned GrafInsight repo to and then run `yarn install --pure-lockfile`.
3. Navigate to `<GrafInsight_DIR>/packages/grafinsight-toolkit` and then run `yarn link`.
4. Navigate to the directory where your plugin code is and then run `npx grafinsight-toolkit plugin:dev --yarnlink`. This adds all dependencies required by grafinsight-toolkit to your project, as well as link your local grafinsight-toolkit version to be used by the plugin.

### Debug grafinsight-toolkit

To debug grafinsight-toolkit you can use standard [NodeJS debugging methods](https://nodejs.org/de/docs/guides/debugging-getting-started/#enable-inspector) (`node --inspect`, `node --inspect-brk`).

To run grafinsight-toolkit in a debugging session use the following command in the toolkit's directory:

`node --inspect-brk ./bin/grafinsight-toolkit.js [task]`

To run [linked](#develop-grafinsight-toolkit) grafinsight-toolkit in a debugging session use the following command in the plugin's directory:

`node --inspect-brk ./node_modules/@grafinsight/toolkit/bin/grafinsight-toolkit.js [task]`
