<template>
  <v-btn
    :icon="icon"
    :size="size"
    :variant="variant"
    :color="color"
    @click="openHelp"
    :title="tooltip"
  ></v-btn>
</template>

<script setup>
const router = useRouter()

const props = defineProps({
  // Path to the help document (e.g., "devices/overview.md")
  docPath: {
    type: String,
    required: true
  },
  // Tooltip text
  tooltip: {
    type: String,
    default: 'View help documentation'
  },
  // Icon to display
  icon: {
    type: String,
    default: 'mdi-help-circle-outline'
  },
  // Button size
  size: {
    type: String,
    default: 'small'
  },
  // Button variant
  variant: {
    type: String,
    default: 'text'
  },
  // Button color
  color: {
    type: String,
    default: 'grey'
  }
})

function openHelp() {
  // Open help in a new window with the document path as a query parameter
  const helpUrl = router.resolve({
    path: '/help',
    query: { doc: props.docPath }
  }).href
  window.open(helpUrl, '_blank', 'width=1200,height=800')
}
</script>
