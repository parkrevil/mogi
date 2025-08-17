export interface NavigationItem {
  icon: string;
  label: string;
  href: string;
}

export const navigation: NavigationItem[] = [
  {
    icon: '🏠',
    label: '홈',
    href: '/'
  },
  {
    icon: '📈',
    label: '영수증',
    href: '/receipts'
  },
  {
    icon: '🏆',
    label: '랭킹',
    href: '/rankings'
  },
  {
    icon: '🎉',
    label: '파티 메이커',
    href: '/party-maker'
  },
];
