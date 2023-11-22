import React from "react"

export interface InputProps
  extends React.InputHTMLAttributes<HTMLInputElement> {
    error?: string,
    label?: string
}

const Input = React.forwardRef<HTMLInputElement, InputProps>(
    ({ className, type, error, label, ...props }, ref) => {
        return (
            <div className="form-control w-full px-2">
                {label ? <label className="label">
                    <span className="label-text">{label}</span>
                </label> : null}
                <input
                    type={type}
                    className="input input-bordered w-full"
                    ref={ref}
                    {...props}
                />
                {error ? <label className="label">
                    <span className="label-text-alt">{error}</span>
                </label> : null}
            </div>
        )
    }
)
Input.displayName = "Input"

export { Input }