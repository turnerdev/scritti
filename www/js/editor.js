import { define, dispatch, html, store } from 'https://unpkg.com/hybrids@latest/src/index.js';
import AssetStore from './store/asset-store.js';

async function selectLines(host) {
  // const source = host.shadowRoot.querySelector('textarea').value.split("\n")
  dispatch(host, 'selectlines', { detail: host.focusedLines.map(i => host.source[i]) });
}

function autofocus(host, event) {
  console.log("AUTOFOCUS", host.focused)
  const textarea = host.shadowRoot.querySelector('textarea')
  textarea && textarea.focus()
}

function updateLine(host, target) {
    setTimeout(() => {
      const start = target.selectionStart;
      const end = target.selectionEnd;
      const startLine = target.value.substr(0, start).split("\n").length;
      const endLine = target.value.substr(0, end).split("\n").length;
      host.focusedLines = Array(endLine - startLine + 1)
        .fill(startLine - 1)
        .map((v, i) => v + i)
    }, 1)
}

function handleMouseDown(host, e) {
  updateLine(host, e.target)
}

function handleMouseMove(host, e) {
  if (e.buttons) {
    updateLine(host, e.target)
  }
}

async function handleKeyDown(host, e) {
  updateLine(host, e.target)

  // Ctrl+0-9: Change focus
  if(e.keyCode >= 48 && e.keyCode <= 58 && e.ctrlKey) {
    dispatch(host, "navigate", { detail: e.keyCode })
  }

  // Enter: New line with current indentation level
  if(e.keyCode == 13) {
    await (async () => {
      const textarea = e.target;
      const value = textarea.value;
      const textLines = value.substr(0, textarea.selectionStart).split("\n");
      const start = e.target.selectionStart;
      const indent = "\n" + textLines[textLines.length-1].match(/(\s)*/)[0]
      textarea.value = textarea.value.slice(0, start) + indent + textarea.value.slice(start)
      textarea.selectionStart = start + indent.length
      textarea.selectionEnd = start + indent.length
      
      updateAsset(host)
      e.preventDefault()
    })()
  }

  // Ctrl+Space: complete suggestion
  if(e.keyCode === 32 && e.ctrlKey && host.currentSuggestion) {
    (() => {
      const textarea = e.target;
      const value = textarea.value;
      const start = e.target.selectionStart;
      textarea.value = value.slice(0, start) + host.currentSuggestion + value.slice(start)
      textarea.selectionStart = start + host.currentSuggestion.length
      textarea.selectionEnd = start + host.currentSuggestion.length

      updateAsset(host)
      e.preventDefault()
    })()
  }

  // Tab, Shift-Tab: Block indentation
  if(e.keyCode === 9) {
    const textarea = e.target;
    const value = textarea.value + "\n";
    const lines = value.split("\n");
    const start = textarea.selectionStart;
    const end = textarea.selectionEnd;
    const startLine = value.substr(0, start).split("\n").length-1;
    const endLine = value.substr(0, end).split("\n").length-1;

    if (e.shiftKey) {
      // Increase indentation for select lines
      const slines = Array.from({length: endLine-startLine+1 }, (_, i) => i + startLine);
      const ltos = value.slice(0, start).split("\n"); // Split value up to selection start into lines
      const cursorAtLineStart = ltos[ltos.length-1].length === 0; // Is the selection at the start of the line?

      const result = slines.map(i => lines[i]).reduce((a, next, i) => {
        const firstSelectionLine = i === 0
        const firstLine = startLine + i === 0;
        if (next.charAt(0) === "\t") {
          a.value += (firstLine ? "" : "\n") + next.slice(1);
          if (firstSelectionLine && !cursorAtLineStart) {
            a.selectionStart -= 1;
          }
          a.selectionEnd -= 1;
        } else {
          a.value += (firstLine ? "" : "\n") + next;
        }

        return a
      }, {
        value: lines.slice(0, startLine).join("\n"),
        selectionStart: textarea.selectionStart,
        selectionEnd: textarea.selectionEnd
      })
      e.target.value = result.value + "\n" + lines.slice(endLine+1).join("\n");
      e.target.value = e.target.value.slice(0, e.target.value.length-1);
      e.target.selectionStart = result.selectionStart;
      e.target.selectionEnd = result.selectionEnd;

    } else {
      // Decrease indentation for selected lines
      const slines = Array.from({length: endLine-startLine+1 }, (_, i) => i + startLine);

      const result = slines.map(i => lines[i]).reduce((a, next, i) => {
        const firstLine = startLine + i === 0;
        a.value += (firstLine ? "\t" : "\n\t") + next;
        if (i === 0) {
          a.selectionStart +=1;
        }
        a.selectionEnd += 1;

        return a
      }, {
        value: lines.slice(0, startLine).join("\n"),
        selectionStart: textarea.selectionStart,
        selectionEnd: textarea.selectionEnd
      })
      e.target.value = result.value + "\n" + lines.slice(endLine+1).join("\n");
      e.target.value = e.target.value.slice(0, e.target.value.length-1);
      e.target.selectionStart = result.selectionStart;
      e.target.selectionEnd = result.selectionEnd;
    }

    e.preventDefault();
    updateAsset(host)
  }
}

function handleKeyUp(host, e) {
  updateLine(host, e.target)
}

async function handleInput(host, e) {
  const textarea = e.target;
  var textLines = textarea.value.substr(0, textarea.selectionStart).split("\n");
  var currentLineNumber = textLines.length;
  var currentColumnIndex = textLines[textLines.length-1].length;
  console.log("Current Line Number "+ currentLineNumber+" Current Column Index "+currentColumnIndex );
  updateAsset(host)
}

/**
 * Updates the Asset Store with the current buffer
 * @param {object} host 
 */
async function updateAsset(host) {
  console.info("Updating Asset source from buffer")
  const textarea = host.shadowRoot.querySelector('textarea');
  host.bufferHistory = [...host.bufferHistory, textarea.value];
  selectLines(host);

  const ready = store.pending(host.asset);
  if (ready) {
    await ready
  }
  store.set(host.asset, {
    source: textarea.value
  });
}

const Editor = {
  assetType: 0,
  assetName: '',
  currentBuffer: {
    get: ({ bufferHistory }) => bufferHistory[bufferHistory.length-1].split('\n')
  },
  bufferHistory: [''],
  focused: false,
  asset: {
    ...store(AssetStore, (host) => `${host.assetType} ${host.assetName}`),
    observe: async (host, value, lastValue) => {
      const asset = await (store.pending(value) || Promise.resolve(value));
      host.source = asset.source.split("\n");
      if (lastValue && asset.id !== lastValue.id && asset.source != host.bufferHistory[host.bufferHistory.length-1]) {
        host.bufferHistory = ['-'];
      }
    }
  },
  focusedLines: {
    connect: (host, key) => {
      host[key] = [];
    },
    observe: (host, value, lastValue) => {
      if (!Array.isArray(value) ||
          !Array.isArray(lastValue) ||
          value.length !== lastValue.length ||
          !value.every((val, i) => val === lastValue[i])) {
        console.log(value);
        selectLines(host);
      }
    }
  },
  currentSuggestion: ({focusedLines, suggestions, source}) => {
    if (focusedLines.length === 1 && source[focusedLines[0]]) {
      const lastToken = source[focusedLines[0]].trimStart();
      const filteredSuggestions = suggestions.filter((suggestion) => suggestion.startsWith(lastToken));
      
      if (filteredSuggestions.length > 1) {
        const next = filteredSuggestions.reduce((a, b) => {
          const maxlength = Math.min(a.length, b.length)
          let i = lastToken.length - 1
          while (++i < maxlength) {
            if (a.charAt(i) != b.charAt(i)) {
              break;
            }
          }
          return a.slice(0, i)
        })
        return next.slice(lastToken.length);
      } else if (filteredSuggestions.length === 1) {
        return filteredSuggestions[0].slice(lastToken.length)
      }
      return '';
    }
    return '';
  },
  suggestions: [],
  source: {
    connect: (host, key) => {
      host[key] = [];
    },
    observe: (host, value, lastValue) => {
      console.info(`Observed source change ${host.assetName}`)

      // Naive method of limiting rerendering of textarea
      // All changes are tracked in a buffer history, if the latest asset source is found in
      // in the buffer history, assume it came from this session and do not update the textarea
      const newSource = value.join("\n");
      const oldSource = lastValue && lastValue.join("\n");

      if (newSource !== oldSource) {
        const bufferIndex = host.bufferHistory.indexOf(newSource);
        if (~bufferIndex && bufferIndex === host.bufferHistory.length - 1) {
          console.log('Inbound source matches current buffer')
        } else if (~bufferIndex) {
          console.info('Inbound source found in history', bufferIndex, host.bufferHistory)
          host.bufferHistory = [...host.bufferHistory.slice(0, bufferIndex), ...host.bufferHistory.slice(bufferIndex+1)];
        } else {
          console.info('Inbound source not found in history, updating editor source', newSource);
          host.shadowRoot.querySelector('textarea').value = newSource;
          host.bufferHistory = [...host.bufferHistory, newSource];
          selectLines(host)
        }
      }
    }
  },
  title: '',
  render: ({ focused, currentBuffer, currentSuggestion, focusedLines, assetName, title }) => html`
    <style>
      :host {
        background-color: #262626;
        color: white;
        display: flex;
        flex-direction: column;
        font-family: 'Lucida Console', Monaco, monospace;
        font-size: 0.9rem;
      }
      .editor {
        display: flex;
        flex-direction: row;
        flex: 1;
        overflow-y: scroll;
        overflow-x: hidden;
        scrollbar-width: thin;
        scrollbar-color: #555 #333
      }
      header {
        background-color: #222222;
        border-color: #161616;
        border-style: solid;
        border-width: 1px 0;
        display: flex;
        flex-direction: row;
        padding: 6px;
      }
      .source {
        flex: 1;
        min-height: 100%;
        position: relative;
        height: min-content;
        word-break: break-all;
        padding-left: 2px;
      }
      .source ul {
        display: block;
        list-style: none;
        margin: 0;
        padding: 0;
        width: 100%;
        position: relative;
      }
      .source ul li {
        -moz-tab-size : 2;
        -o-tab-size : 2;
        color: rgba(255,255,255, 0.4);
        counter-increment: item;
        display: flex;
        margin: 0;
        padding: 0;
        position: relative;
        tab-size : 2;
      }
      .source ul li:before {
        content: counter(item);
        counter-rest: none;
        left: -62px;
        position: absolute;
        text-align: right;
        width: 50px;
      }
      .sourceline {
        white-space: pre-wrap;
        max-width: 100%;
      }
      .error {
        background-color: #550011;
      }
      .focus {
        background-color: #333333;
      }
      .wrapper {
        display: flex;
        flex-direction: row;
        width: 100%;
      }
      textarea {
        -moz-tab-size : 2;
        -o-tab-size : 2;
        background: transparent;
        border: none;
        box-sizing: border-box;
        color: #AAA;
        display: block;
        flex: 1;
        font-family: 'Lucida Console', Monaco, monospace;
        font-size: 0.9rem;
        height: 100%;
        margin: 0;
        overflow: hidden;
        padding: 0 2px 0 0;
        position: absolute;
        resize: none;
        tab-size : 2;
        top: 0;
        width: 100%;
      }
      textarea::-moz-selection { 
        background: #115577;
      }
      textarea::selection {
        background: #115577;
      }
      textarea:focus {
        color: #FFF;
        outline: none;
      }
      .linenumbers {
        background-color: #222222;
        border-right: 1px solid #161616;
        color: #AAA;
        color: transparent;
        line-height: normal;
        padding: 0 10px;
        text-align: right;
        user-select: none;
      }
      .title {
        color: #666;
        margin-left: auto;
      }
    </style>
    <header class=${{focused}}>
      <div class="asset-name">${assetName}</div>
      <div class="title">${title}</div>
    </header>
    <div class="editor">
      <div class="wrapper">
        <div class="linenumbers">
          ${currentBuffer.map((_, lineNo) => html`${lineNo+1}<br>`)}
        </div>
        <div class="source">
          <ul>
            ${currentBuffer.map((sourceline, i) => html`
              <li class="${{'focus':  !!~focusedLines.indexOf(i)}}">
                <div class="sourceline">${sourceline}</div>
                ${!!~focusedLines.indexOf(i) && html`<div class="autocomplete">${currentSuggestion}</div>`}
              </li>
            `)}
          </ul>
          <textarea
            spellcheck="false"
            onmousedown="${handleMouseDown}"
            onmousemove="${handleMouseMove}"
            onkeydown="${handleKeyDown}"
            onkeyup="${handleKeyUp}"
            oninput="${handleInput}"
            >${focused && autofocus}</textarea>
        </div>
      </div>
    </div>
  `,
};

define('x-editor', Editor);

export default Editor;