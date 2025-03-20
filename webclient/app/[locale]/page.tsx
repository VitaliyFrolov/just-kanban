import { FC } from 'react';
import { getLocale } from 'next-intl/server';

import { RoutePath } from 'shared/config';
import { redirect } from 'shared/lib';

const RoutingPage: FC = async () => {
  const locale = await getLocale();
  return redirect({ href: RoutePath.Home, locale });
};

export default RoutingPage;
