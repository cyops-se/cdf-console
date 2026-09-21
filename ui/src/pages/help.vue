<route lang="yaml">
meta:
  layout: minimal
</route>

<template>
  <v-container fluid class="help-container pa-0">
    <v-row no-gutters class="help-row">
      <!-- Left sidebar with folder structure -->
      <div v-if="showSidebar" class="doc-sidebar">
        <v-card class="fill-height" elevation="0">
          <v-card-text class="pa-0">
            <v-progress-linear v-if="loading" indeterminate color="primary"></v-progress-linear>
            <v-list v-if="docStructure" density="compact">
              <template v-for="folder in docStructure.children" :key="folder.path">
                <v-list-group v-if="folder.isDir" :value="folder.name">
                  <template v-slot:activator="{ props }">
                    <v-list-item v-bind="props" :prepend-icon="'mdi-folder'">
                      <v-list-item-title>{{ folder.title || folder.name }}</v-list-item-title>
                    </v-list-item>
                  </template>
                  <v-list-item
                    v-for="file in folder.children"
                    :key="file.path"
                    :value="file.path"
                    @click="loadDocContent(file.path)"
                    :active="activePath === file.path"
                    :prepend-icon="'mdi-file-document'"
                    class="pl-2"
                  >
                    <v-list-item-title>{{ file.title || file.name }}</v-list-item-title>
                  </v-list-item>
                </v-list-group>
                <v-list-item
                  v-else
                  :value="folder.path"
                  @click="loadDocContent(folder.path)"
                  :active="activePath === folder.path"
                  :prepend-icon="'mdi-file-document'"
                >
                  <v-list-item-title>{{ folder.title || folder.name }}</v-list-item-title>
                </v-list-item>
              </template>
            </v-list>
          </v-card-text>
          <div v-if="sysinfo.gitversion || sysinfo.gitcommit" class="pa-3 text-center version-info">
            <div v-if="sysinfo.gitversion" class="text-caption text-grey">
              Version: {{ sysinfo.gitversion }}
            </div>
            <div v-if="sysinfo.gitcommit" class="text-caption text-grey">
              Commit: {{ sysinfo.gitcommit.substring(0, 7) }}
            </div>
          </div>
        </v-card>
      </div>

      <!-- Right content area -->
      <div class="doc-content">
        <v-card class="fill-height" elevation="0">
          <v-card-text class="pa-6 doc-content-scroll">
            <v-alert v-if="error" type="error" class="mb-4">
              {{ error }}
            </v-alert>
            <v-progress-circular v-if="loadingContent" indeterminate color="primary"></v-progress-circular>
            <div v-else-if="currentContent" ref="pdfContentRef" :class="markdownClass" @click="handleLinkClick" v-html="renderedHtml"></div>
            <div v-else class="text-center text-grey">
              <v-icon size="64" class="mb-4">mdi-book-open-variant</v-icon>
              <p>Välj något från menyn till vänster!</p>
            </div>
          </v-card-text>
        </v-card>
      </div>
    </v-row>
  </v-container>
</template>

<script setup>
import { ref, onMounted, watch, computed, inject } from 'vue'
import { useRouter } from 'vue-router'
import { useTheme } from 'vuetify'
import axios from 'axios'
import { marked } from 'marked'

// Define route meta if needed for auto-routing
defineOptions({
  name: 'Help'
})

const router = useRouter()
const theme = useTheme()
const showSidebar = inject('showSidebar', ref(true))
const activePath = inject('activePath', ref(''))
const pdfContent = inject('pdfContent', ref(null))
const sysinfo = inject('sysinfo', ref({ gitversion: '', gitcommit: '' }))
const pdfContentRef = ref(null)
const docStructure = ref(null)
const currentContent = ref('')
const loading = ref(false)
const loadingContent = ref(false)
const error = ref('')

// Compute markdown class based on theme
const markdownClass = computed(() => {
  return theme.global.current.value.dark ? 'markdown-body markdown-dark' : 'markdown-body markdown-light'
})

// Update pdfContent when the ref changes
watch(pdfContentRef, (newVal) => {
  pdfContent.value = newVal
})

// Process markdown content to convert image paths to API URLs and render to HTML
const renderedHtml = computed(() => {
  if (!currentContent.value) return ''

  // Replace relative image paths with API endpoint URLs
  // Matches: ![alt text](images/filename.png) or ![](images/file.jpg)
  const processedMarkdown = currentContent.value.replace(
    /!\[([^\]]*)\]\(images\/([^)]+)\)/g,
    '![$1](/api/help/images/$2)'
  )

  // Convert markdown to HTML using marked
  return marked(processedMarkdown)
})

// Load the documentation structure
async function loadDocStructure() {
  loading.value = true
  error.value = ''
  try {
    const response = await axios.get('/api/help/structure')
    if (response.data.success) {
      docStructure.value = response.data.root
      // Auto-load README.md if it exists
      if (docStructure.value.children) {
        const readme = docStructure.value.children.find(c => c.name === '00-overview.md')
        if (readme) {
          await loadDocContent(readme.path)
          activePath.value = readme.path
        }
      }
    } else {
      error.value = response.data.errors.join(', ')
    }
  } catch (err) {
    error.value = 'Failed to load documentation structure: ' + err.message
  } finally {
    loading.value = false
  }
}

// Load the content of a specific document
async function loadDocContent(path) {
  if (!path || path === 'docs') return // Don't load folder content

  loadingContent.value = true
  error.value = ''
  try {
    const response = await axios.get(`/api/help/content/${path}`)
    if (response.data.success) {
      currentContent.value = response.data.content
      activePath.value = path
    } else {
      error.value = response.data.errors.join(', ')
      currentContent.value = ''
    }
  } catch (err) {
    error.value = 'Failed to load document: ' + err.message
    currentContent.value = ''
  } finally {
    loadingContent.value = false
  }
}

// Handle clicks on links in markdown content
function handleLinkClick(event) {
  const target = event.target
  if (target.tagName === 'A') {
    const href = target.getAttribute('href')
    // Handle relative markdown links
    if (href && href.endsWith('.md') && !href.startsWith('http')) {
      event.preventDefault()
      // Resolve relative path based on current document
      const currentDir = activePath.value.substring(0, activePath.value.lastIndexOf('/'))
      let newPath = href

      if (href.startsWith('../')) {
        // Go up one directory
        const parentDir = currentDir.substring(0, currentDir.lastIndexOf('/'))
        newPath = parentDir + '/' + href.substring(3)
      } else if (!href.startsWith('/')) {
        // Same directory
        newPath = currentDir + '/' + href
      } else {
        // Absolute path (remove leading slash)
        newPath = href.substring(1)
      }

      loadDocContent(newPath)
    }
  }
}

// Watch for changes in query parameters
watch(() => router.currentRoute.value.query.doc, (newDoc) => {
  if (newDoc) {
    loadDocContent(newDoc)
  }
})

onMounted(() => {
  loadDocStructure()

  // Check if a document was specified in the query parameter
  const docParam = router.currentRoute.value.query.doc
  if (docParam) {
    // Wait a bit for the structure to load, then load the specified document
    setTimeout(() => {
      loadDocContent(docParam)
    }, 500)
  }
})
</script>

<style scoped>
.help-container {
  height: 100%;
  max-height: 100%;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.help-row {
  height: 100%;
  max-height: 100%;
  overflow: hidden;
  flex: 1;
  display: flex;
  flex-direction: row;
}

.doc-sidebar {
  width: 300px;
  min-width: 300px;
  border-right: 1px solid rgba(0, 0, 0, 0.12);
  overflow: hidden;
  height: 100%;
}

.doc-sidebar :deep(.v-card) {
  display: flex;
  flex-direction: column;
}

.doc-sidebar :deep(.v-card-text) {
  overflow-y: auto;
  flex: 1;
}

.doc-sidebar .version-info {
  border-top: 1px solid rgba(0, 0, 0, 0.12);
  line-height: 1.4;
}

.doc-sidebar .version-info .text-caption {
  font-size: 0.7rem;
}

.doc-content {
  flex: 1;
  height: 100%;
  overflow: hidden;
}

.doc-content :deep(.v-card) {
  display: flex;
  flex-direction: column;
}

.doc-content-scroll {
  overflow-y: auto;
  flex: 1;
}

/* Smaller font size for menu items */
.doc-sidebar :deep(.v-list-item-title) {
  font-size: 0.875rem;
}

.doc-sidebar :deep(.v-list-group__header .v-list-item-title) {
  font-size: 0.875rem;
  font-weight: 500;
}

/* Reduce space between icon and text */
.doc-sidebar :deep(.v-list-item__prepend) {
  width: 24px !important;
  min-width: 24px !important;
  margin-inline-end: 0 !important;
}

.doc-sidebar :deep(.v-list-item__spacer) {
  width: 4px !important;
}

.doc-sidebar :deep(.v-list-item) {
  column-gap: 4px !important;
}

/* Reduce subitem indentation */
.doc-sidebar :deep(.v-list-group__items .v-list-item) {
  padding-inline-start: 50px !important;
}

/* GitHub markdown styling */
.markdown-body {
  box-sizing: border-box;
  width: 100%;
  padding: 10px;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Helvetica, Arial, sans-serif;
  font-size: 16px;
  line-height: 1.6;
  word-wrap: break-word;
  orphans: 3;
  widows: 3;
}

/* Light theme colors */
.markdown-light {
  color: #24292f;
  background-color: #ffffff;
}

.markdown-light :deep(h1),
.markdown-light :deep(h2),
.markdown-light :deep(h3),
.markdown-light :deep(h4),
.markdown-light :deep(h5),
.markdown-light :deep(h6) {
  margin-top: 10px;
  margin-bottom: 5px;
  font-weight: 600;
  line-height: 1.25;
  color: #24292f;
  page-break-after: avoid;
  break-after: avoid;
  page-break-inside: avoid;
  break-inside: avoid;
}

.markdown-light :deep(h1) { font-size: 2em; border-bottom: 1px solid #d0d7de; padding-bottom: 0.3em; }
.markdown-light :deep(h2) { font-size: 1.5em; border-bottom: 1px solid #d0d7de; padding-bottom: 0.3em; }
.markdown-light :deep(h3) { font-size: 1.25em; }
.markdown-light :deep(h4) { font-size: 1em; }
.markdown-light :deep(h5) { font-size: 0.875em; }
.markdown-light :deep(h6) { font-size: 0.85em; color: #57606a; }

.markdown-light :deep(p) {
  margin-top: 0;
  margin-bottom: 16px;
  page-break-inside: avoid !important;
  break-inside: avoid !important;
  orphans: 3;
  widows: 3;
}

.markdown-light :deep(a) { color: #0969da; text-decoration: none; }
.markdown-light :deep(a:hover) { text-decoration: underline; }

.markdown-light :deep(code) {
  padding: 0.2em 0.4em;
  margin: 0;
  font-size: 85%;
  background-color: rgba(175, 184, 193, 0.2);
  border-radius: 6px;
  font-family: ui-monospace, SFMono-Regular, 'SF Mono', Menlo, Consolas, 'Liberation Mono', monospace;
}

.markdown-light :deep(pre) {
  padding: 16px;
  overflow: auto;
  font-size: 85%;
  line-height: 1.45;
  background-color: #f6f8fa;
  border-radius: 6px;
  margin-bottom: 16px;
  page-break-inside: avoid;
  break-inside: avoid;
}

.markdown-light :deep(ul),
.markdown-light :deep(ol) {
  padding-left: 2em;
  margin-top: 0;
  margin-bottom: 16px;
  page-break-inside: avoid;
  break-inside: avoid;
}

.markdown-light :deep(li) {
  page-break-inside: avoid;
  break-inside: avoid;
}

.markdown-light :deep(img) {
  max-width: 100%;
  box-sizing: border-box;
  background-color: #ffffff;
  page-break-inside: avoid !important;
  break-inside: avoid !important;
  page-break-before: auto;
  page-break-after: auto;
  display: block;
  margin: 16px 0;
}

.markdown-light :deep(p img) {
  display: inline-block;
  vertical-align: middle;
}

/* Dark theme colors */
.markdown-dark {
  color: #c9d1d9;
  background-color: transparent;
}

.markdown-dark :deep(h1),
.markdown-dark :deep(h2),
.markdown-dark :deep(h3),
.markdown-dark :deep(h4),
.markdown-dark :deep(h5),
.markdown-dark :deep(h6) {
  margin-top: 10px;
  margin-bottom: 5px;
  font-weight: 600;
  line-height: 1.25;
  color: #c9d1d9;
  page-break-after: avoid;
  break-after: avoid;
  page-break-inside: avoid;
  break-inside: avoid;
}

.markdown-dark :deep(h1) { font-size: 2em; border-bottom: 1px solid #21262d; padding-bottom: 0.3em; }
.markdown-dark :deep(h2) { font-size: 1.5em; border-bottom: 1px solid #21262d; padding-bottom: 0.3em; }
.markdown-dark :deep(h3) { font-size: 1.25em; }
.markdown-dark :deep(h4) { font-size: 1em; }
.markdown-dark :deep(h5) { font-size: 0.875em; }
.markdown-dark :deep(h6) { font-size: 0.85em; color: #8b949e; }

.markdown-dark :deep(p) {
  margin-top: 0;
  margin-bottom: 16px;
  page-break-inside: avoid !important;
  break-inside: avoid !important;
  orphans: 3;
  widows: 3;
}

.markdown-dark :deep(a) { color: #58a6ff; text-decoration: none; }
.markdown-dark :deep(a:hover) { text-decoration: underline; }

.markdown-dark :deep(code) {
  padding: 0.2em 0.4em;
  margin: 0;
  font-size: 85%;
  background-color: rgba(110, 118, 129, 0.4);
  border-radius: 6px;
  font-family: ui-monospace, SFMono-Regular, 'SF Mono', Menlo, Consolas, 'Liberation Mono', monospace;
}

.markdown-dark :deep(pre) {
  padding: 16px;
  overflow: auto;
  font-size: 85%;
  line-height: 1.45;
  background-color: #161b22;
  border-radius: 6px;
  margin-bottom: 16px;
  page-break-inside: avoid;
  break-inside: avoid;
}

.markdown-dark :deep(ul),
.markdown-dark :deep(ol) {
  padding-left: 2em;
  margin-top: 0;
  margin-bottom: 16px;
  page-break-inside: avoid;
  break-inside: avoid;
}

.markdown-dark :deep(li) {
  page-break-inside: avoid;
  break-inside: avoid;
}

.markdown-dark :deep(img) {
  max-width: 100%;
  box-sizing: border-box;
  background-color: transparent;
  page-break-inside: avoid !important;
  break-inside: avoid !important;
  page-break-before: auto;
  page-break-after: auto;
  display: block;
  margin: 16px 0;
}

.markdown-dark :deep(p img) {
  display: inline-block;
  vertical-align: middle;
}

@media (max-width: 767px) {
  .markdown-body {
    padding: 15px;
  }
}
</style>
