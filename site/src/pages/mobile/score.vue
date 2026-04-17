<template>
  <div class="mobile-score-page">
    <MobileHeader title="积分中心" :show-back="true" />
    
    <div class="score-content">
      <!-- 积分卡片 - 简化版 -->
      <div class="score-card">
        <div class="score-header">
          <div class="score-info">
            <div class="score-label">我的积分</div>
            <div class="score-value">{{ user?.score || 0 }}</div>
          </div>
          <div class="checkin-info" v-if="checkIn">
            <div class="checkin-label">连续签到</div>
            <div class="checkin-days">{{ checkIn.consecutiveDays || 0 }} 天</div>
          </div>
        </div>
      </div>

      <!-- 积分总排行 -->
      <div class="rank-section">
        <div class="rank-title">积分排行</div>
        <div v-if="scoreRank && scoreRank.length > 0" class="rank-list">
          <div v-for="(rankUser, index) in scoreRank" :key="rankUser.id" class="rank-item">
            <div class="rank-number" :class="{ 'top3': index < 3 }">
              <span v-if="index < 3">{{ ['🥇', '🥈', '🥉'][index] }}</span>
              <span v-else>{{ index + 1 }}</span>
            </div>
            <MyAvatar :user="rankUser" :size="32" class="rank-avatar" />
            <div class="rank-info">
              <div class="rank-name">{{ rankUser.nickname }}</div>
              <div class="rank-desc">{{ rankUser.topicCount }} 话题 · {{ rankUser.commentCount }} 评论</div>
            </div>
            <div class="rank-score">
              <i class="iconfont icon-score" />
              <span>{{ rankUser.score }}</span>
            </div>
          </div>
        </div>
        <div v-else class="empty-rank">
          暂无排行数据
        </div>
      </div>

      <!-- 今日签到排行 -->
      <div class="rank-section checkin-rank">
        <div class="rank-title">今日签到排行</div>
        <div v-if="checkInRank && checkInRank.length > 0" class="rank-list">
          <div v-for="(rank, index) in checkInRank" :key="rank.id" class="rank-item">
            <div class="rank-number" :class="{ 'top3': index < 3 }">
              <span v-if="index < 3">{{ ['🥇', '🥈', '🥉'][index] }}</span>
              <span v-else>{{ index + 1 }}</span>
            </div>
            <MyAvatar :user="rank.user" :size="32" class="rank-avatar" />
            <div class="rank-info">
              <div class="rank-name">{{ rank.user.nickname }}</div>
              <div class="rank-time">@{{ usePrettyDate(rank.updateTime) }}</div>
            </div>
          </div>
        </div>
        <div v-else class="empty-rank">
          暂无签到记录
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
const userStore = useUserStore()
const { user } = storeToRefs(userStore)

// 签到数据
const { data: checkIn } = await useMyFetch(`/api/checkin/checkin`)

// 今日签到排行
const { data: checkInRank } = await useMyFetch(`/api/checkin/rank`)

// 积分总排行
const { data: scoreRank } = await useMyFetch(`/api/user/score/rank`)
</script>

<style scoped>
.mobile-score-page {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.score-content {
  padding: 12px;
}

/* 积分卡片 - 紧凑版 */
.score-card {
  background: white;
  border-radius: 10px;
  padding: 16px;
  margin-bottom: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.score-header {
  display: flex;
  justify-content: space-around;
  align-items: center;
}

.score-info, .checkin-info {
  text-align: center;
}

.score-label, .checkin-label {
  font-size: 12px;
  color: #666;
  margin-bottom: 4px;
}

.score-value {
  font-size: 28px;
  font-weight: 700;
  color: #f59e0b;
}

.checkin-days {
  font-size: 20px;
  font-weight: 600;
  color: #3273dc;
}

/* 排行区域 - 紧凑版 */
.rank-section {
  background: white;
  border-radius: 10px;
  padding: 12px;
  margin-bottom: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.checkin-rank {
  margin-bottom: 0;
}

.rank-title {
  font-size: 14px;
  font-weight: 600;
  color: #333;
  margin-bottom: 10px;
  padding-bottom: 8px;
  border-bottom: 1px solid #eee;
}

.rank-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.rank-item {
  display: flex;
  align-items: center;
  padding: 8px;
  background: #f8f8f8;
  border-radius: 6px;
}

.rank-number {
  width: 28px;
  font-size: 16px;
  text-align: center;
}

.rank-number.top3 {
  font-size: 18px;
}

.rank-avatar {
  margin: 0 8px;
}

.rank-info {
  flex: 1;
  min-width: 0;
}

.rank-name {
  font-size: 13px;
  font-weight: 600;
  color: #333;
  margin-bottom: 2px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.rank-desc, .rank-time {
  font-size: 11px;
  color: #999;
}

.rank-score {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #f59e0b;
  font-weight: 600;
  background: #fef3c7;
  padding: 4px 8px;
  border-radius: 12px;
}

.rank-score i {
  font-size: 12px;
}

.empty-rank {
  text-align: center;
  padding: 24px 20px;
  color: #999;
  font-size: 13px;
}
</style>
