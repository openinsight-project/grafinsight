import { feedToDataFrame } from './utils';
import { RssFeed, NewsItem } from './types';
import { DataFrameView } from '@grafinsight/data';

describe('news', () => {
  test('convert RssFeed to DataFrame', () => {
    const frame = feedToDataFrame(grafinsight20191216);
    expect(frame.length).toBe(5);

    // Iterate the links
    const view = new DataFrameView<NewsItem>(frame);
    const links = view.map((item: NewsItem) => {
      return item.link;
    });
    expect(links).toEqual([
      'https://grafinsight.com/blog/2019/12/13/meet-the-grafinsight-labs-team-aengus-rooney/',
      'https://grafinsight.com/blog/2019/12/12/register-now-grafinsightcon-2020-is-coming-to-amsterdam-may-13-14/',
      'https://grafinsight.com/blog/2019/12/10/pro-tips-dashboard-navigation-using-links/',
      'https://grafinsight.com/blog/2019/12/09/how-to-do-automatic-annotations-with-grafinsight-and-loki/',
      'https://grafinsight.com/blog/2019/12/06/meet-the-grafinsight-labs-team-ward-bekker/',
    ]);
  });
});

const grafinsight20191216 = {
  items: [
    {
      title: 'Meet the GrafInsight Labs Team: Aengus Rooney',
      link: 'https://grafinsight.com/blog/2019/12/13/meet-the-grafinsight-labs-team-aengus-rooney/',
      pubDate: 'Fri, 13 Dec 2019 00:00:00 +0000',
      content: '\n\n<p>As GrafInsight Labs continues to grow, we&rsquo;d like you to get to know the team members...',
    },
    {
      title: 'Register Now! GrafInsightCon 2020 Is Coming to Amsterdam May 13-14',
      link: 'https://grafinsight.com/blog/2019/12/12/register-now-grafinsightcon-2020-is-coming-to-amsterdam-may-13-14/',
      pubDate: 'Thu, 12 Dec 2019 00:00:00 +0000',
      content: '\n\n<p>Amsterdam, we&rsquo;re coming back!</p>\n\n<p>Mark your calendars for May 13-14, 2020....',
    },
    {
      title: 'Pro Tips: Dashboard Navigation Using Links',
      link: 'https://grafinsight.com/blog/2019/12/10/pro-tips-dashboard-navigation-using-links/',
      pubDate: 'Tue, 10 Dec 2019 00:00:00 +0000',
      content:
        '\n\n<p>Great dashboards answer a limited set of related questions. If you try to answer too many questions in a single dashboard, it can become overly complex. ...',
    },
    {
      title: 'How to Do Automatic Annotations with GrafInsight and Loki',
      link: 'https://grafinsight.com/blog/2019/12/09/how-to-do-automatic-annotations-with-grafinsight-and-loki/',
      pubDate: 'Mon, 09 Dec 2019 00:00:00 +0000',
      content:
        '\n\n<p>GrafInsight annotations are great! They clearly mark the occurrence of an event to help operators and devs correlate events with metrics. You may not be aware of this, but GrafInsight can automatically annotate graphs by ...',
    },
    {
      title: 'Meet the GrafInsight Labs Team: Ward Bekker',
      link: 'https://grafinsight.com/blog/2019/12/06/meet-the-grafinsight-labs-team-ward-bekker/',
      pubDate: 'Fri, 06 Dec 2019 00:00:00 +0000',
      content:
        '\n\n<p>As GrafInsight Labs continues to grow, we&rsquo;d like you to get to know the team members who are building the cool stuff you&rsquo;re using. Check out the latest of our Friday team profiles.</p>\n\n<h2 id="meet-ward">Meet Ward!</h2>\n\n<p><strong>Name:</strong> Ward...',
    },
  ],
  feedUrl: 'https://grafinsight.com/blog/index.xml',
  title: 'Blog on GrafInsight Labs',
  description: 'Recent content in Blog on GrafInsight Labs',
  generator: 'Hugo -- gohugo.io',
  link: 'https://grafinsight.com/blog/',
  language: 'en-us',
  lastBuildDate: 'Fri, 13 Dec 2019 00:00:00 +0000',
} as RssFeed;
