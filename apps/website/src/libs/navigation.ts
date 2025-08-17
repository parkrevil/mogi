export interface NavigationItem {
  label: string;
  href: string;
}

export const navigation: NavigationItem[] = [
  {
    label: '홈',
    href: '/'
  },
  {
    label: '영수증',
    href: '/receipts'
  },
  {
    label: '랭킹',
    href: '/rankings'
  },
  {
    label: '파티 메이커',
    href: '/party-maker'
  },
];
