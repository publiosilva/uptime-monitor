import LinkButton from '../../components/LinkButton'
import Logo from '../../components/Logo'
import Input from '../../components/Input'
import Button from '../../components/Button'
import { useState } from 'react'
import { z } from 'zod'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useNavigate } from 'react-router-dom'
import { authService } from '../../services/auth'
import { ApiError } from '../../lib/api'

const signUpSchema = z.object({
    email: z
        .email('Invalid email address')
        .min(1, 'Email is required'),
    password: z
        .string()
        .min(1, 'Password is required')
        .min(8, 'Password must be at least 8 characters'),
    confirmPassword: z
        .string()
        .min(1, 'Password is required')
        .min(8, 'Password must be at least 8 characters'),
})

type SignUpForm = z.infer<typeof signUpSchema>

export default function SignUp() {
    const navigate = useNavigate()
    const [submitError, setSubmitError] = useState<string | null>(null)
    const {
        register,
        handleSubmit,
        formState: { errors, isSubmitting }
    } = useForm<SignUpForm>({
        resolver: zodResolver(signUpSchema),
        defaultValues: {
            email: '',
            password: '',
            confirmPassword: '',
        }
    })

    async function onSubmit(data: SignUpForm) {
        if (data.password !== data.confirmPassword) {
            setSubmitError('Passwords do not match')
        }

        try {
            await authService.register(data)
            navigate('/')
        } catch (error) {
            if (error instanceof ApiError) {
                setSubmitError(error.message)
                return
            }

            setSubmitError('Something went wrong. Please try again.')
        }
    }

    return (
        <div className="flex flex-col items-center justify-center h-screen">
            <div className="z-[-1] absolute inset-0 opacity-20" style={{
                backgroundImage: 'linear-gradient(rgb(26, 45, 74) 1px, transparent 1px), linear-gradient(90deg, rgb(26, 45, 74) 1px, transparent 1px)',
                backgroundSize: '40px 40px',
                backgroundPosition: '0 0, 0 0',
                backgroundRepeat: 'repeat',
            }}></div>
            <div className="flex flex-col items-center justify-center h-screen gap-6 w-full max-w-sm">
                <div className="flex items-center gap-2">
                    <Logo />
                    <span className="text-2xl font-semibold text-slate-100 tracking-tight">uptime</span>
                </div>
                <div className="flex flex-col bg-[#0c1422] border border-[#1a2d4a] rounded-2xl p-6 shadow-2xl w-full">
                    <h1 className="text-lg font-semibold text-slate-100 mb-0.5">Create your account</h1>
                    <p className="text-xs text-[#3a5070] mb-6">Start monitoring in under 2 minutes</p>
                    <form className="flex flex-col gap-4" onSubmit={handleSubmit(onSubmit)} noValidate>
                        {submitError && (
                            <p className="text-xs text-red-400 bg-red-500/10 border border-red-500/20 rounded-lg px-3 py-2">
                                {submitError}
                            </p>
                        )}
                        <Input
                            id="email"
                            label="Email address"
                            type="email"
                            placeholder="example@email.com"
                            autoComplete="email"
                            error={errors.email?.message}
                            {...register('email')}
                        />
                        <Input
                            id="password"
                            label="Password"
                            type="password"
                            placeholder="••••••••"
                            error={errors.password?.message}
                            {...register('password')}
                        />
                        <Input
                            id="confirm-password"
                            label="Confirm password"
                            type="password"
                            placeholder="••••••••"
                            error={errors.confirmPassword?.message}
                            {...register('confirmPassword')}
                        />
                        <Button type="submit" disabled={isSubmitting}>
                            {isSubmitting ? 'Creating account…' : 'Create account'}
                        </Button>
                    </form>
                    <p className="text-xs text-[#3a5070] text-center mt-5">Already have an account? <LinkButton to="/signin">Sign in</LinkButton></p>
                </div>
                <p className="text-center text-xs text-[#2a3e58] mt-6">
                    Self-hosted · Open source · No telemetry
                </p>
            </div>
        </div>
    )
}
