# GrafInsight UI components library

> **@grafinsight/ui is currently in BETA**.

@grafinsight/ui is a collection of components used by [GrafInsight](https://github.com/grafinsight/grafinsight)

Our goal is to deliver GrafInsight's common UI elements for plugins developers and contributors.

See [package source](https://github.com/grafinsight/grafinsight/tree/master/packages/grafinsight-ui) for more details.

## Installation

`yarn add @grafinsight/ui`

`npm install @grafinsight/ui`

## Development

For development purposes we suggest using `yarn link` that will create symlink to @grafinsight/ui lib. To do so navigate to `packages/grafinsight-ui` and run `yarn link`. Then, navigate to your project and run `yarn link @grafinsight/ui` to use the linked version of the lib. To unlink follow the same procedure, but use `yarn unlink` instead.

### Storybook 6.x migration

We've upgraded Storybook to version 6 and with that we will convert to using [controls](https://storybook.js.org/docs/react/essentials/controls) instead of knobs for manipulating components. Controls will not require as much coding as knobs do. Please refer to the [storybook style-guide](https://github.com/grafinsight/grafinsight/blob/master/contribute/style-guides/storybook.md#contrls) for further information.
