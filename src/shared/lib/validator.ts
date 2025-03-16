import { ZodIssueCode, z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';

export {
  z as validate,
  zodResolver as formResolver,
  ZodIssueCode as IssueCode,
};
