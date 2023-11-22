import React from "react"
import { Control, Controller, FieldValue } from "react-hook-form"
import { Input, InputProps } from "./input"

export interface InputControllerProps extends InputProps {
    control: any
    defaultValue: string
    rules: any
    name: string
}

const InputController = ({ control, type, rules, defaultValue, name, ...props }: InputControllerProps) => {
    return (
        <Controller
            control={control}
            name={name}
            rules={{ required: "Required" }}
            render={({ field: { onChange, onBlur, value, ref }, fieldState: { error } }) => (
                <Input onChange={onChange} onBlur={onBlur} value={value} ref={ref} error={error?.message} {...props} />
            )}
        />
    )
}

export default InputController