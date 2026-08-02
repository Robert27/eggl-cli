import {
  HomeBackground,
  HomeFeature,
  HomeFooter,
  HomeHero,
  Layout as BasicLayout,
  renderHtmlOrText,
  type HomeLayoutProps,
} from "@rspress/core/theme-original";
import { useSite } from "@rspress/core/runtime";
import { HeroEgglGraphic } from "./HeroEgglGraphic";

function SiteFooterMessage() {
  const { site } = useSite();
  const message = site.themeConfig.footer?.message;
  if (!message) {
    return null;
  }

  return (
    <div className="eggl-site-footer">
      <div
        className="eggl-site-footer__message"
        {...renderHtmlOrText(message)}
      />
    </div>
  );
}

function HomeLayout(props: HomeLayoutProps) {
  const {
    beforeHero,
    afterHero,
    beforeFeatures,
    afterFeatures,
    beforeHeroActions,
    afterHeroActions,
  } = props;

  return (
    <>
      <HomeBackground />
      {beforeHero}
      <HomeHero
        beforeHeroActions={beforeHeroActions}
        afterHeroActions={afterHeroActions}
        image={<HeroEgglGraphic />}
      />
      {afterHero}
      {beforeFeatures}
      <HomeFeature />
      {afterFeatures}
      <HomeFooter />
    </>
  );
}

const Layout = () => <BasicLayout afterDocFooter={<SiteFooterMessage />} />;

export * from "@rspress/core/theme-original";
export { Layout, HomeLayout };
