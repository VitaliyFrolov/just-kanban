import { createNavigation } from 'next-intl/navigation';

import { LOCALE_PREFIX, SUPPORTED_LOCALES } from 'shared/config/intl';

export const { Link, usePathname, useRouter, redirect, getPathname } =
  createNavigation({
    locales: SUPPORTED_LOCALES,
    localePrefix: LOCALE_PREFIX,
  });
