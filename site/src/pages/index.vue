<template>
  <section class="main">
    <div class="container main-container" style="display: flex; gap: 15px;">
      <div class="left-container" style="width: 260px; min-width: 260px; max-width: 260px;">
        <check-in />
        <site-notice />
        <score-rank />
        <topics-nav style="margin-top: 20px;" />
      </div>
      <div class="right-container" style="flex: 1; min-width: 0;">
        <div class="topics-main">
          <load-more-async v-slot="{ results }" url="/api/topic/topics">
            <topic-list :topics="results" show-sticky />
          </load-more-async>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
useHead({
  title: useSiteTitle(),
});
</script>

<style lang="scss" scoped></style>