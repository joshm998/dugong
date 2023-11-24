import type { InputProps } from "./input";
import { Input } from "./input"
import { useField } from "remix-validated-form";

export interface InputFieldProps extends InputProps {
    name: string,
    label: string
}

const InputField = ({ name, label }: InputFieldProps) => {
    const { error, getInputProps } = useField(name);
    return (
      <Input {...getInputProps({ id: name })} error={error} label={label} />
    );
  };

export default InputField