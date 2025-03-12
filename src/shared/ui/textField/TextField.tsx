import { FC, Ref, useId } from 'react';

import { FormControl, InformationText, Input, InputProps } from 'shared/ui';
import { clsx } from 'shared/lib';

export interface TextFieldProps extends InputProps {
  label?: string;
  helperText?: string;
  multiline?: boolean;
  ref?: Ref<HTMLInputElement>;
}

export const TextField: FC<TextFieldProps> = ({
  className,
  helperText,
  label,
  multiline,
  ref,
  id: idProp,
  ...otherProps
}) => {
  const innerId = useId();
  const id = idProp ?? innerId;

  return (
    <div className={className}>
      <FormControl {...otherProps}>
        {label && (
          <label
            className="text-secondary font-medium mb-[5px] ml-[2px] cursor-pointer"
            htmlFor={id}
          >
            {label}
          </label>
        )}
        {multiline ? (
          <input
            {...otherProps}
            className={clsx(className, 'field-sizing-content')}
            ref={ref}
            id={id}
          />
        ) : (
          <Input
            {...otherProps}
            ref={ref}
            id={id}
          />
        )}
        {helperText && (
          <InformationText
            className="mt-1"
            error={otherProps.error}
          >
            {helperText}
          </InformationText>
        )}
      </FormControl>
    </div>
  );
};

TextField.displayName = 'TextField';
