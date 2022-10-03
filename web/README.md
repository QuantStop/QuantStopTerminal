<h1 align="center">QuantStopTerminal</h1>
<div align="center">

[![Build Status][build-status-img]][build-status-url]
[![License][license-img]][license-url]
[![GoDoc][godoc-img]][godoc-url]
[![Go Report Card][go-report-img]][go-report-url]
[![Twitter URL][twitter-img]][twitter-url]

</div>
<br>
<div align="center">

[![Logo][logo-img]][logo-url]

</div>
<br>
<p align="center">
  An open source trading and analysis platform.
  <br>
  <a href="https://github.com/quantstop/quantstopterminal"><strong>Explore the docs »</strong></a>
  <br>
  <br>
  <a href="https://github.com/quantstop/quantstopterminal">View Demo</a>
  ·
  <a href="https://github.com/QuantStop/QuantStopTerminal/issues/new?assignees=&labels=&template=bug_report.md&title=">Report Bug</a>
  ·
  <a href="https://github.com/QuantStop/QuantStopTerminal/issues/new?assignees=&labels=&template=feature_request.md&title=">Request Feature</a>
</p>

## Recommended IDE Setup

[VSCode](https://code.visualstudio.com/) + [Volar](https://marketplace.visualstudio.com/items?itemName=Vue.volar) (and disable Vetur) + [TypeScript Vue Plugin (Volar)](https://marketplace.visualstudio.com/items?itemName=Vue.vscode-typescript-vue-plugin).

## Type Support for `.vue` Imports in TS

TypeScript cannot handle type information for `.vue` imports by default, so we replace the `tsc` CLI with `vue-tsc` for type checking. In editors, we need [TypeScript Vue Plugin (Volar)](https://marketplace.visualstudio.com/items?itemName=Vue.vscode-typescript-vue-plugin) to make the TypeScript language service aware of `.vue` types.

If the standalone TypeScript plugin doesn't feel fast enough to you, Volar has also implemented a [Take Over Mode](https://github.com/johnsoncodehk/volar/discussions/471#discussioncomment-1361669) that is more performant. You can enable it by the following steps:

1. Disable the built-in TypeScript Extension
    1) Run `Extensions: Show Built-in Extensions` from VSCode's command palette
    2) Find `TypeScript and JavaScript Language Features`, right click and select `Disable (Workspace)`
2. Reload the VSCode window by running `Developer: Reload Window` from the command palette.

## Customize configuration

See [Vite Configuration Reference](https://vitejs.dev/config/).

## Project Setup

```sh
npm install
```

### Compile and Hot-Reload for Development

```sh
npm run dev
```

### Type-Check, Compile and Minify for Production

```sh
npm run build
```

### Run Unit Tests with [Vitest](https://vitest.dev/)

```sh
npm run test:unit
```

### Run End-to-End Tests with [Cypress](https://www.cypress.io/)

```sh
npm run build
npm run test:e2e # or `npm run test:e2e:ci` for headless testing
```

### Lint with [ESLint](https://eslint.org/)

```sh
npm run lint
```



<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[build-status-url]: https://github.com/quantstop/quantstopterminal/actions/workflows/release.yml/badge.svg?branch=release
[build-status-img]: https://github.com/quantstop/quantstopterminal/actions/workflows/release.yml/badge.svg?branch=release
[license-url]: https://github.com/quantstop/quantstopterminal/blob/release/LICENSE
[license-img]: https://img.shields.io/badge/License-MIT-orange.svg?style=flat-round
[godoc-url]: https://godoc.org/github.com/quantstop/quantstopterminal
[godoc-img]: https://godoc.org/github.com/quantstop/quantstopterminal?status.svg
[go-report-url]: https://goreportcard.com/report/github.com/quantstop/quantstopterminal
[go-report-img]: https://goreportcard.com/badge/github.com/quantstop/quantstopterminal
[twitter-url]: https://twitter.com/quantstop
[twitter-img]: https://img.shields.io/badge/twitter-@QuantStop-wnZunKusqrz0QZNxE4Ag?logo=twitter&style=flat
[logo-url]: https://quantstop.com
[logo-img]: ../assets/images/qst.png
[vue-url]: https://vuejs.org/
[vue-img]: ../assets/images/vuejs-logo.png

