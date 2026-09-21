<template>
  <v-layout class="help-layout">
    <v-app-bar color="primary" density="compact">
      <v-btn icon="mdi-menu" @click="toggleSidebar" title="Toggle menu"></v-btn>
      <v-app-bar-title class="text-h6">Help Documentation</v-app-bar-title>
      <template v-slot:append>
        <v-btn v-if="activePath" icon="mdi-file-pdf-box" @click="downloadPDF" title="Download as PDF"></v-btn>
        <v-btn icon="mdi-theme-light-dark" @click="toggleTheme" title="Toggle theme"></v-btn>
      </template>
    </v-app-bar>

    <v-main class="help-main">
      <router-view />
    </v-main>
  </v-layout>
</template>

<style scoped>
.help-layout {
  height: 100vh;
  overflow: hidden;
}

.help-main {
  height: 100%;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}
</style>

<script setup>
import { ref, provide, onMounted } from 'vue'
import { useTheme } from 'vuetify'
import html2pdf from 'html2pdf.js'
import axios from 'axios'

const theme = useTheme()
const showSidebar = ref(true)
const activePath = ref('')
const pdfContent = ref(null)
const sysinfo = ref({
  gitversion: '',
  gitcommit: ''
})

provide('showSidebar', showSidebar)
provide('activePath', activePath)
provide('pdfContent', pdfContent)
provide('sysinfo', sysinfo)

onMounted(async () => {
  try {
    const response = await axios.get('/api/system/sysinfo')
    if (response.data.success && response.data.sysinfo) {
      sysinfo.value = response.data.sysinfo
    }
  } catch (error) {
    console.error('Failed to fetch system info:', error)
  }
})

function toggleTheme() {
  theme.global.name.value = theme.global.current.value.dark ? 'light' : 'dark'
}

function toggleSidebar() {
  showSidebar.value = !showSidebar.value
}

async function downloadPDF() {
  if (!pdfContent.value) return

  // Get the filename from the active path
  const filename = activePath.value.split('/').pop().replace('.md', '.pdf')

  // Clone the content element
  const element = pdfContent.value.cloneNode(true)

  // Force light theme for PDF
  element.classList.remove('markdown-dark')
  element.classList.add('markdown-light')

  // Configure html2pdf options
  const opt = {
    margin: 15,
    filename: filename,
    image: { type: 'jpeg', quality: 0.98 },
    html2canvas: {
      scale: 2,
      useCORS: true,
      letterRendering: true,
      allowTaint: true
    },
    jsPDF: { unit: 'mm', format: 'a4', orientation: 'portrait' },
    pagebreak: {
      mode: ['avoid-all', 'css', 'legacy'],
      before: '.page-break-before',
      after: '.page-break-after',
      avoid: ['p', 'h1', 'h2', 'h3', 'h4', 'h5', 'h6', 'img', 'pre', 'ul', 'ol', 'li']
    }
  }

  // Generate PDF
  await html2pdf().set(opt).from(element).save()
}
</script>
