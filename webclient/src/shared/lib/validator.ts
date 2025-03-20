import { ZodIssueCode, z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';

export {
  z as validate,
  zodResolver as formResolver,
  ZodIssueCode as IssueCode,
};

export interface ApiValidationError<
  TFields extends Record<string, string> = Record<string, string>,
> {
  message: string;
  fields: TFields;
}
