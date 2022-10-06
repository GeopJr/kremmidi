<h1 align="center">onion-rings</h1>
<h4 align="center">Webpage that compiles multiple parts of a binary into the original file.</h4>

# Building

```bash
pnpm i
pnpm dev # for dev mode
pnpm build # to build it
```

After building it, it should be a **single** `index.html` under `dist/`.

# Screenshots

<p align="center">
    <img width="512" src="https://i.imgur.com/MtvK7ZS.png" alt="
        A webpage with the following content:
        Drag and drop all parts here.
        Please make sure to upload ALL parts.
        -- Drag and drop component with the text: Drag 'n' drop some files here, or click to select files --
        -- Save button with the label: Save --
        Files:
         1. tor-browser-11.5.3-android-armv7-multi.apk.000
         2. tor-browser-11.5.3-android-armv7-multi.apk.001
         3. tor-browser-11.5.3-android-armv7-multi.apk.002
         4. tor-browser-11.5.3-android-armv7-multi.apk.003
         5. tor-browser-11.5.3-android-armv7-multi.apk.004
         6. tor-browser-11.5.3-android-armv7-multi.apk.005
         7. tor-browser-11.5.3-android-armv7-multi.apk.006
         8. tor-browser-11.5.3-android-armv7-multi.apk.007
         9. tor-browser-11.5.3-android-armv7-multi.apk.008
        10. tor-browser-11.5.3-android-armv7-multi.apk.009
        11. tor-browser-11.5.3-android-armv7-multi.apk.010
        This website will combine all parts into one binary. If the extension is missing, please add it yourself.
        If any part is missing, the file will be incomplete.
        If you uploaded the wrong file, please refresh the page.
    "/>
</p>

## FaQ

- Why pnpm?

Both npm and yarn should still work. Pnpm uses symlinks for node_modules instead of downloading them each time for all your projects.

- Why Svelte?

Mostly reactivity.
