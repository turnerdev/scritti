import {
  define,
  html,
  store,
} from "https://unpkg.com/hybrids@latest/src/index.js"
import AssetStore from "./store/asset-store.js"
import ApplicationStore from "./store/application-store.js"
import Editor from "./editor.js"
import Preview from "./preview.js"
import Source from "./source.js"

function handleNavigate(host, event) {
  host.focusEditor = event.detail
}

function handleComponentSelectLines(host, event) {
  if (event.detail.length > 0) {
    const line = event.detail[0].trim()
    const regex = /^(?<indent>\s*)((?<tag>\w+)\.)?(?<style>\w+)?(\s*"(?<text>[^"\\]*(\.[^"\\]*)*)")?\s*$/
    const m = regex.exec(line)
    host.style = m?.groups?.style || ''
    host.svg = m?.groups?.tag === 'svg'
  }
}

function handleClasses(host, event) {
  host.styleSuggestions = event.detail.sort()
}

function toggleEditor(host) {
  if (window.location.hash === '#editor') {
    window.history.replaceState(null, null, ' ');
  } else {
    window.history.replaceState(null, null, '#editor');
  }
  window.dispatchEvent(new HashChangeEvent('hashchange'));
}

function toggleViewSource(host) {
  host.sourceVisible = !host.sourceVisible
}

const IDE = {
  app: store(ApplicationStore),
  style: "",
  svg: false,
  styleSuggestions: [],
  editorVisible: {
    connect: (host, key) => {
      host[key] = true;
      window.onhashchange = (e) => {
        host[key] = window.location.hash === '#editor'
      }
    }
  },
  sourceVisible: false,
  focusEditor: 49,
  stylesheet: 'https://unpkg.com/tailwindcss@^2/dist/tailwind.min.css',
  render: ({ app, focusEditor, style, stylesheet, styleSuggestions, svg, editorVisible, sourceVisible }) => html`
    <style>
      :host {
        display: flex;
        height: 100%;
        width: 100%;
        font-family: Seravek, Calibri, Roboto, Arial, sans-serif;
      }
      .editor {
        background-color: gray;
        display: flex;
        flex-direction: column;
        width: 40vw;
        margin-left: -40vw;
        transition: margin 0.6s ease;
      }
      .editor.visible {
        margin: 0;
      }
      x-editor,
      x-preview {
        flex: 1 1 0;
        height: 0;
      }
      .preview {
        flex: 2;
        position: relative;
        overflow-y: scroll;
      }
      header {
        align-items: center;
        background-color: #111;
        display: flex;
        flex-direction: row;
        justify-content: space-between;
        padding: 6px;
      }
      header .app-name {
        color: white;
        font-weight: bold;
      }
      header .app-status > div {
        border-radius: 8px;
        color: white;
        display: none;
        font-size: 0.75rem;
        font-weight: bold;
        padding: 4px 8px;
        text-transform: uppercase;
      }
      header .app-status .connected {
        display: block;
        background-color: #0a751b;
      }
      header .app-status .local {
        display: block;
        background-color: #ca3900;
      }
      .stylesheet {
        background-color: #111;
        color: #888;
        display: flex;
      }
      .stylesheet label {
        font-size: 1rem;
        margin-right: 6px;
        padding: 6px;
      }
      .stylesheet input {
        background-color: #444;
        border-radius: 0;
        border: none;
        color: white;
        flex: 1;
        padding: 4px 8px;
      }
      .edit {
        background-color: rgba(0,0,0,0.2);
        color: white;
        cursor: pointer;
        height: 20px;
        padding: 7px;
        position: fixed;
        top: 0;
        width: 20px;
        z-index: 100;
      }
      .viewsource {
        align-items: center;
        background-color: rgba(0,0,0,0.3);
        color: white;
        cursor: pointer;
        display: flex;
        font-size: 0.9rem;
        height: 20px;
        padding: 7px;
        position: fixed;
        right: 0;
        top: 0;
        white-space: nowrap;
        width: min-content;
        z-index: 100;
      }
      .viewsource svg {
        width: 20px;
        height: 20px;
        padding-right: 5px;
      }
      .viewsource:hover {
        background-color: rgba(0,0,0,0.5);
      }
    </style>
    ${app.ready &&
    html`
      <div class="${{editor: true, visible: editorVisible}}">
        <header>
          <div class="app-name">scritti</div>
          <div class="app-status">
            <div
              class="${{
                connected: app.ready && app.server,
                local: app.ready && !app.server,
              }}"
            >
              ${app.server ? "connected" : "local"}
            </div>
          </div>
        </header>
        <div class="stylesheet">
          <label>stylesheet</label>
          <input value="${stylesheet}" onchange="${html.set('stylesheet')}"/>
        </div>
        <x-editor
          asset-type="0"
          asset-name="main"
          title="ctrl+1"
          focused="${focusEditor === 49}"
          onnavigate="${handleNavigate}"
          onselectlines="${handleComponentSelectLines}"
        >
        </x-editor>
        ${style &&
        html`
          <x-editor
            asset-type="1"
            assetName="${style}"
            title="ctrl+2"
            focused="${focusEditor === 50}"
            onnavigate="${handleNavigate}"
            suggestions="${styleSuggestions}"
          >
          </x-editor>
        `}
        ${svg &&
        html`
          <x-editor
            asset-type="2"
            assetName="${style}"
            title="ctrl+3"
            focused="${focusEditor === 51}"
            onnavigate="${handleNavigate}"
          >
          </x-editor>
        `}
      </div>
      <div class="preview">
        <div class="edit" onclick="${toggleEditor}">
          ${!editorVisible ? html`
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>`:html`
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>`
        }
        </div>
        <div class="viewsource" onclick="${toggleViewSource}">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M12.316 3.051a1 1 0 01.633 1.265l-4 12a1 1 0 11-1.898-.632l4-12a1 1 0 011.265-.633zM5.707 6.293a1 1 0 010 1.414L3.414 10l2.293 2.293a1 1 0 11-1.414 1.414l-3-3a1 1 0 010-1.414l3-3a1 1 0 011.414 0zm8.586 0a1 1 0 011.414 0l3 3a1 1 0 010 1.414l-3 3a1 1 0 11-1.414-1.414L16.586 10l-2.293-2.293a1 1 0 010-1.414z" clip-rule="evenodd" />
          </svg>
          ${sourceVisible ? 'Hide Source' : 'Show Source'}
        </div>
        ${sourceVisible
          ? html`<x-source
            asset-type="0"
            asset-name="main"></x-source>`
          : html`<x-preview
            asset-type="0"
            asset-name="main"
            stylesheet="${stylesheet}"
            onclasses="${handleClasses}"></x-preview>`
        }
      </div>
    `}
  `,
}

define("x-ide", IDE)
