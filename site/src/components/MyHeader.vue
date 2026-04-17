<template>
  <!-- 移动端导航栏 -->
  <nav
    v-if="needMobile"
    class="mobile-navbar"
    role="navigation"
    aria-label="main navigation"
  >
    <div class="mobile-nav-content">
      <nuxt-link :to="logoLink" class="mobile-logo">
        <img
          v-if="config.siteLogo"
          :alt="config.siteTitle"
          :src="config.siteLogo"
          class="logo-img"
        />
        <img v-else :alt="config.siteTitle" src="~/assets/images/logo.png" class="logo-img" />
      </nuxt-link>
      <span class="mobile-site-title">{{ config.siteTitle }}</span>
      <div class="mobile-spacer"></div>
      <div class="mobile-create-btn" @click="handleCreateTopic">
        <i class="iconfont icon-plus"></i>
      </div>
      <a
        :class="{ 'is-active': navbarActive }"
        class="navbar-burger burger"
        data-target="navbarBasic"
        @click="toggleNav"
      >
        <span aria-hidden="true" />
        <span aria-hidden="true" />
        <span aria-hidden="true" />
      </a>
    </div>
    <div :class="{ 'is-active': navbarActive }" class="navbar-menu mobile-menu">
      <div class="navbar-start">
        <nuxt-link
          v-for="(nav, index) in filteredMobileNavs"
          :key="index"
          :to="nav.url"
          class="navbar-item"
          @click="navbarActive = false"
        >
          {{ nav.title }}
        </nuxt-link>
      </div>

      <div class="navbar-end">
        <!-- 积分和签到 -->
        <div v-if="isLogin" class="mobile-score-section">
          <nuxt-link class="navbar-item" to="/mobile/score" @click="navbarActive = false">
            <i class="iconfont icon-gold" />
            <span>积分: {{ user.score || 0 }}</span>
          </nuxt-link>
          <a class="navbar-item" @click="handleMobileCheckIn">
            <i class="iconfont icon-calendar" />
            <span v-if="checkIn && checkIn.checkIn">已连续签到 {{ checkIn.consecutiveDays }} 天</span>
            <span v-else>签到</span>
          </a>
        </div>

        <msg-notice v-if="user" />

        <div
          v-if="user"
          class="navbar-item has-dropdown is-hoverable user-menus"
        >
          <div class="navbar-link">
            <MyAvatar :user="user" :size="24" />
            <span
              :to="`/user/${user.id}`"
              class="user-menus-nickname ellipsis"
              >{{ user.nickname }}</span
            >
          </div>
          <div class="navbar-dropdown">
            <nuxt-link class="navbar-item" :to="`/user/${user.id}`" @click="navbarActive = false">
              <i class="iconfont icon-username" />
              <span>{{ $t("common.header.profile") }}</span>
            </nuxt-link>
            <nuxt-link class="navbar-item" to="/user/favorites" @click="navbarActive = false">
              <i class="iconfont icon-favorite" />
              <span>{{ $t("common.header.favorites") }}</span>
            </nuxt-link>
            <nuxt-link class="navbar-item" to="/user/profile" @click="navbarActive = false">
              <i class="iconfont icon-edit" />
              <span>{{ $t("common.header.editProfile") }}</span>
            </nuxt-link>
            <a class="navbar-item" @click="signout">
              <i class="iconfont icon-log-out" />
              <span>{{ $t("common.header.logout") }}</span>
            </a>
          </div>
        </div>
        <div v-else class="navbar-item">
          <div class="buttons">
            <nuxt-link class="button login-btn" to="/user/signin">
              {{ $t("common.header.login") }}
            </nuxt-link>
          </div>
        </div>
      </div>
    </div>
  </nav>

  <!-- 桌面端导航栏 -->
  <nav
    v-else
    class="navbar has-shadow is-fixed-top"
    role="navigation"
    aria-label="main navigation"
  >
    <div class="container">
      <div class="navbar-brand">
        <nuxt-link :to="logoLink" class="navbar-item">
          <img
            v-if="config.siteLogo"
            :alt="config.siteTitle"
            :src="config.siteLogo"
          />
          <img v-else :alt="config.siteTitle" src="~/assets/images/logo.png" />
        </nuxt-link>
        <a
          :class="{ 'is-active': navbarActive }"
          class="navbar-burger burger"
          data-target="navbarBasic"
          @click="toggleNav"
        >
          <span aria-hidden="true" />
          <span aria-hidden="true" />
          <span aria-hidden="true" />
        </a>
      </div>
      <div :class="{ 'is-active': navbarActive }" class="navbar-menu">
        <div class="navbar-start">
          <nuxt-link
            v-for="(nav, index) in config.siteNavs"
            :key="index"
            :to="nav.url"
            class="navbar-item"
          >
            {{ nav.title }}
          </nuxt-link>
        </div>

        <div class="navbar-end">
          <div class="navbar-item">
            <search-input />
          </div>

          <div class="navbar-item">
            <create-topic-btn />
          </div>

          <msg-notice v-if="user" />

          <div
            v-if="user"
            class="navbar-item has-dropdown is-hoverable user-menus"
          >
            <div class="navbar-link">
              <MyAvatar :user="user" :size="24" />
              <span
                :to="`/user/${user.id}`"
                class="user-menus-nickname ellipsis"
                >{{ user.nickname }}</span
              >
            </div>
            <div class="navbar-dropdown">
              <nuxt-link class="navbar-item" :to="`/user/${user.id}`">
                <i class="iconfont icon-username" />
                <span>{{ $t("common.header.profile") }}</span>
              </nuxt-link>
              <nuxt-link class="navbar-item" to="/user/favorites">
                <i class="iconfont icon-favorite" />
                <span>{{ $t("common.header.favorites") }}</span>
              </nuxt-link>
              <nuxt-link class="navbar-item" to="/user/profile">
                <i class="iconfont icon-edit" />
                <span>{{ $t("common.header.editProfile") }}</span>
              </nuxt-link>
              <a class="navbar-item" @click="signout">
                <i class="iconfont icon-log-out" />
                <span>{{ $t("common.header.logout") }}</span>
              </a>
            </div>
          </div>
          <div v-else class="navbar-item">
            <div class="buttons">
              <nuxt-link class="button login-btn" to="/user/signin">
                {{ $t("common.header.login") }}
              </nuxt-link>
            </div>
          </div>
          <div class="navbar-item">
            <color-mode />
          </div>
        </div>
      </div>
    </div>
  </nav>
</template>

<script setup>
const userStore = useUserStore();
const configStore = useConfigStore();
const route = useRoute();
const { isWeCom, isMobile } = useWeCom();

const { user } = storeToRefs(userStore);
const { config } = storeToRefs(configStore);
const { t } = useI18n();

const navbarActive = ref(false);

// 判断是否需要移动端布局
const needMobile = computed(() => {
  return isWeCom.value || isMobile.value
})

// 移动端过滤后的导航（去掉"话题"和"文章"）
const filteredMobileNavs = computed(() => {
  if (!config.value.siteNavs) return []
  return config.value.siteNavs.filter(nav => {
    // 过滤掉标题为"话题"或"文章"的导航项
    return nav.title !== '话题' && nav.title !== '文章'
  })
})

// Logo 链接
const logoLink = computed(() => {
  return needMobile.value ? '/mobile' : '/'
})

// 登录状态
const isLogin = computed(() => {
  return userStore.user !== null;
});

// 签到数据
const { data: checkIn, refresh: refreshCheckIn } = await useMyFetch(
  `/api/checkin/checkin`,
  { server: false, immediate: isLogin.value }
);

function toggleNav() {
  navbarActive.value = !navbarActive.value;
}

async function signout() {
  if (confirm(t("common.header.confirmLogout"))) {
    await userStore.signout();
    useLinkTo("/");
  }
}

// 移动端签到
async function handleMobileCheckIn() {
  if (!isLogin.value) {
    useLinkTo('/user/signin');
    return;
  }
  if (checkIn.value && checkIn.value.checkIn) {
    return; // 已签到
  }
  try {
    await useHttpPost("/api/checkin/checkin");
    useMsgSuccess(t("component.checkIn.checkInSuccess"));
    await refreshCheckIn();
  } catch (e) {
    useCatchError(e);
  }
}

function handleCreateTopic() {
  const router = useRouter();
  router.push('/topic/create');
}
</script>

<style lang="scss" scoped>
.navbar {
  /*opacity: 0.99;*/
  /*border-bottom: 1px solid #e7edf3;*/
  background-color: var(--bg-color);

  .navbar-item {
    font-weight: 700;
  }

  .publish {
    color: var(--text-color);
    background-color: #3174dc;
    width: 100px;
    &:hover {
      color: var(--text-color);
      background-color: #4d91fa;
    }
  }

  .login-btn {
    height: 32px;
    border-color: #000; // TODO
    &:hover {
      color: var(--text-color3);
      border-color: var(--text-color3);
    }
  }
}

.user-menus {
  .navbar-link {
    display: flex;
    align-items: center;

    .user-menus-nickname {
      margin-left: 5px;
      padding: 0 4px;
      font-size: 14px;
      color: var(--text-color);

      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
  }
  .navbar-dropdown {
    border: 1px solid var(--border-color);

    a {
      display: flex;
      align-items: center;
      // padding: 8px 16px;
      img {
        width: 20px;
        height: 20px;
      }
      span {
        margin-left: 10px;
        width: 56px;
        height: 20px;
        font-size: 14px;
        font-weight: 400;
        line-height: 20px;
      }
    }
  }
}

// 移动端导航栏样式
.mobile-navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 44px;
  background: white;
  border-bottom: 1px solid #eee;
  z-index: 100;
}

.mobile-nav-content {
  display: flex;
  align-items: center;
  height: 100%;
  padding: 0 12px;
}

.mobile-logo {
  display: flex;
  align-items: center;
  cursor: pointer;
}

.logo-img {
  height: 28px;
  width: auto;
  object-fit: contain;
}

.mobile-site-title {
  margin-left: 8px;
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.mobile-spacer {
  flex: 1;
}

.mobile-create-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  background: #3273dc;
  border-radius: 50%;
  color: white;
  cursor: pointer;
  margin-right: 12px;
}

.mobile-create-btn .iconfont {
  font-size: 16px;
}

.mobile-menu {
  background: white;
  border-bottom: 1px solid #eee;
}
</style>
