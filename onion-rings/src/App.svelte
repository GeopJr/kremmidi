<script>
  import Dropzone from "svelte-file-dropzone";
  import { saveAs } from "file-saver";

  /**
   * @type {{accepted: File[]}} The drag-n-dropped files.
   * @const
   */
  const files = {
    accepted: [],
  };

  function handleFilesSelect(e) {
    const { acceptedFiles } = e.detail;

    files.accepted = [...files.accepted, ...acceptedFiles];
  }

  /**
   * Extensions of accepted
   * files.
   * @type {string[]}
   * @const
   */
  const allowedExtensions = ["tar.xz", "dmg", "exe", "apk"];

  /**
   * Returns the mimetype base on filename.
   * @param {string} filename The filename.
   * @returns {string} The mime type.
   */
  function mimeType(filename) {
    const filenameLower = filename.toLowerCase();

    let mime = "application/x-binary";

    if (filenameLower.endsWith(".apk")) {
      mime = "application/vnd.android.package-archive";
    } else if (filenameLower.endsWith(".tar.xz")) {
      mime = "application/x-xz";
    } else if (filenameLower.endsWith(".exe")) {
      mime = "application/x-msdos-program";
    } else if (filenameLower.endsWith(".dmg")) {
      mime = "application/x-apple-diskimage";
    }

    return mime;
  }

  // Remove anything that doesn't match the regex.
  // If a filename matches regex but has an additional suffix
  // e.g. "filename.000 (1)", it will remove the suffix.

  /**
   * Regex that matches
   * tor-browser-{anything}.000
   * @type {RegExp}
   * @const
   */
  const filenameRegex = /^tor-?browser-.+\.[0-9]{3}/i;

  /**
   * Returns the part's number.
   * @param {string} str The part filename.
   * @returns {number} The part number.
   */
  function getPartNo(str) {
    return parseInt(str.match(filenameRegex)[0].split(".").pop(), 10);
  }

  /**
   * Whether the filename has
   * and allowed extension.
   * @param {string} str The filename.
   * @returns {boolean} Whether it has an allowed extension.
   */
  function hasAllowedExtension(str) {
    const withoutPartNo = removeNumberSuffix(
      str.match(filenameRegex)[0].split(".").slice(0, -1).join(".")
    ).toLowerCase();

    return allowedExtensions.some((s) => withoutPartNo.endsWith(s));
  }

  /**
   * Removes suffixes added from
   * the browser to avoid conflicts
   * e.g. "test (1).png"
   * => "test.png"
   * @param {string} str The string to remove the suffix from.
   * @returns {string} The string without the suffix.
   */
  function removeNumberSuffix(str) {
    return str.replace(/ ?\([0-9]+\) ?/, "");
  }

  $: parts = [...files.accepted].filter(
    (x) => filenameRegex.test(x.name) && hasAllowedExtension(x.name)
  );
  // Sort them based on the .XXX.
  $: sortedParts = parts.sort((a, b) => getPartNo(a.name) - getPartNo(b.name));
  // Whether the save button is disabled.
  $: saveBtnDisabled = sortedParts.length === 0;

  /**
   * Checks if each sorted part number matches its array
   * position.
   * If e.g. i=5, [5] = 6 it means that one of the
   * first five parts is missing.
   * (e.g [0, 1, 2, 3, 4, 6])
   * @param {File[]} cleanParts Array of *sorted* parts.
   * @returns {boolean} Whether it passed or not.
   */
  function smokeTestPartNumbers(cleanParts) {
    const arrOfPartNumbers = cleanParts.map((x) => getPartNo(x.name));

    for (let i = 0; i < arrOfPartNumbers.length; i++) {
      if (i !== arrOfPartNumbers[i]) {
        return false;
      }
    }

    return true;
  }

  /**
   * Checks if all parts have the same extension.
   * @param {File[]} cleanParts Array of *sorted* parts.
   * @returns {boolean} Whether it passed or not.
   */
  function smokeTestPartExtensions(cleanParts) {
    // WARN: ".tar.xz" => "xz" but it has already passed the other
    //       extension filters so it should be fine.
    const exts = cleanParts.map((x) =>
      removeNumberSuffix(
        x.name.match(filenameRegex)[0].split(".").slice(0, -1).pop()
      )
    );

    return new Set(exts).size === 1;
  }

  /**
   * Runs smoke tests.
   * @param {File[]} cleanParts Array of *sorted* parts.
   */
  function runSmokeTests(cleanParts) {
    if (!smokeTestPartNumbers(cleanParts)) {
      alert(
        "Uh oh!\nLooks like one of the parts is missing!\nMake sure you downloaded all of them and keep their original names."
      );
      return window.location.reload();
    }

    if (!smokeTestPartExtensions(cleanParts)) {
      alert(
        "Uh oh!\nLooks like one of the parts has a different extension!\nMake sure you select all files with the same extension and original names."
      );
      return window.location.reload();
    }
  }

  /**
   * Combines the parts into a single file
   * and prompts the browser to save it.
   */
  function save() {
    if (sortedParts.length === 0) return;

    runSmokeTests(sortedParts);

    // The save file name should be
    // the original file (aka without the .XXX).
    const name = removeNumberSuffix(
      `${sortedParts[0].name}`.split(".").slice(0, -1).join(".")
    );

    // Add them back together.
    const blob = new Blob(sortedParts, {
      type: mimeType(name),
    });

    // Ask browser to save.
    saveAs(blob, name);
  }
</script>

<h1>Drag and drop all parts here.</h1>
<h2>Please make sure to upload <b>ALL</b> parts.</h2>
<Dropzone on:drop={handleFilesSelect} />
<button disabled={saveBtnDisabled} on:click={() => save()}>Save</button>
<ol>
  {#if sortedParts.length > 0}
    <h2>Files:</h2>
    {#each sortedParts as item}
      <li><code>{item.name}</code></li>
    {/each}
  {/if}
</ol>

<p>
  This website will combine all parts into one binary. If the extension is
  missing, please add it yourself.
</p>
<p>If any part is missing, the file will be incomplete.</p>
<p>If you uploaded the wrong file, please refresh the page.</p>

<style>
  :root {
    font-family: Inter, Avenir, Helvetica, Arial, sans-serif;
    font-size: 16px;
    line-height: 24px;
    font-weight: 400;

    color-scheme: light dark;
    color: rgba(255, 255, 255, 0.87);
    background-color: #242424;

    font-synthesis: none;
    text-rendering: optimizeLegibility;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    -webkit-text-size-adjust: 100%;
  }

  ol {
    padding: 0;
  }

  h1 {
    font-size: 3.2em;
    line-height: 1.1;
  }

  button {
    border-radius: 8px;
    border: 2px solid transparent;
    padding: 0.6em 1.2em;
    margin: 1rem 0rem;
    font-size: 1em;
    font-weight: 500;
    font-family: inherit;
    background-color: #1a1a1a;
    cursor: not-allowed;
    transition: border-color 0.25s;
  }

  button:enabled {
    cursor: pointer;
  }

  button:hover:enabled {
    border-color: #646cff;
  }

  button:focus:enabled,
  button:focus-visible:enabled {
    outline: 4px auto -webkit-focus-ring-color;
  }

  :global(.dropzone) {
    border-radius: 8px !important;
    border-style: solid !important;
    background-color: #1a1a1a !important;
    border-color: #1a1a1a !important;
    color: #a5a5a5 !important;
  }

  :global(.dropzone):focus-visible {
    outline: auto !important;
  }

  :global(.dropzone):hover {
    border-color: #646cff !important;
  }

  @media (prefers-color-scheme: light) {
    button {
      background-color: #ececec;
    }

    :global(.dropzone) {
      background-color: #ececec !important;
      border-color: #ececec !important;
      color: #474747 !important;
    }
  }
</style>
