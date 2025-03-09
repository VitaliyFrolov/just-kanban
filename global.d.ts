import en from './public/messages/en.json';

type Messages = typeof en;

declare global {
  type IntlMessages = Messages;
}
