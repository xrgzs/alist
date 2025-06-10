# Alist

🗂️A file list program that supports multiple storages, powered by Gin and Solidjs.

This is a fork of https://github.com/AlistGo/alist. And the original project is suspected to have been sold since 2025.

## Document

Been Hacked. **Don't trust any command or link in `alistgo.com` and `alist666` without auditing.**

## Usage

This fork is for development and PR, and only provides the CI version of docker images on `ghcr.io`. There is no image on Docker Hub. And you should not use this if you don't trust me.

```bash
docker pull ghcr.io/xrgzs/alist:main
```

If you need to run it on other platforms, please compile it yourself.

## Compile

1. Install `git`, `go`. And configure GCC following [C/C++ for Visual Studio Code](https://code.visualstudio.com/docs/languages/cpp).

2. Clone `alist` back-end source code.

   ```bash
   git clone https://github.com/xrgzs/alist --depth=1
   ```

3. Build `alist-web` front-end. You can use the pre-builded dist and extract it to `public/dist`:

   ```bash
   curl -L https://codeload.github.com/xrgzs/alist-web/tar.gz/refs/heads/web-dist -o alist-web-web-dist.tar.gz
   tar -zxvf alist-web-web-dist.tar.gz
   rm -rf public/dist
   mv -f alist-web-web-dist/dist public
   rm -rf alist-web-web-dist alist-web-web-dist.tar.gz
   ```

   Or build it yourself. You should install `nodejs` and `pnpm`.

   ```bash
   git clone https://github.com/xrgzs/alist-web --depth=1
   cd alist-web
   pnpm install
   pnpm i18n:build
   pnpm build
   cd ..
   rm -rf public/dist
   cp -r alist-web/dist public/dist
   ```

4. Build binary. Do not use build.sh.

   ```bash
   go build main.go
   ```

## Demo

None.

## Discussion

Disabled.

**Issues are for bug reports and feature requests only.**

## License

The `AList` is open-source software licensed under the AGPL-3.0 license.

## Disclaimer

- This program is a free and open source project. It is designed to share files on the network disk, which is convenient for downloading and learning Golang. Please abide by relevant laws and regulations when using it, and do not abuse it;
- This program is implemented by calling the official sdk/interface, without destroying the official interface behavior;
- This program only does 302 redirect/traffic forwarding, and does not intercept, store, or tamper with any user data;
- Before using this program, you should understand and bear the corresponding risks, including but not limited to account ban, download speed limit, etc., which is none of this program's business;
- If there is any infringement, please contact me by GitHub, and it will be dealt with in time.
