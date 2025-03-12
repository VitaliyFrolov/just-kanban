import 'next-intl';

import { SUPPORTED_LOCALES } from 'shared/config/intl';
import auth from 'public/messages/en/auth.json';
import common from 'public/messages/en/common.json';

declare module 'next-intl' {
  interface AppConfig {
    Messages: {
      common: typeof common;
      auth: typeof auth;
    };
    Locale: (typeof SUPPORTED_LOCALES)[number];
  }
}
