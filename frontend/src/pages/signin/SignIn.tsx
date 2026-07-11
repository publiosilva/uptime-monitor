import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import Logo from '../../components/Logo'
import Input from '../../components/Input'
import Button from '../../components/Button'
import LinkButton from '../../components/LinkButton'
import { ApiError } from '../../lib/api'
import { authService } from '../../services/auth'

const signInSchema = z.object({
    email: z
        .email('Invalid email address')
        .min(1, 'Email is required'),
    password: z
        .string()
        .min(1, 'Password is required')
        .min(8, 'Password must be at least 8 characters'),
})

type SignInForm = z.infer<typeof signInSchema>

export default function SignIn() {
    const navigate = useNavigate()
    const [submitError, setSubmitError] = useState<string | null>(null)

    const {
        register,
        handleSubmit,
        formState: { errors, isSubmitting },
    } = useForm<SignInForm>({
        resolver: zodResolver(signInSchema),
        defaultValues: {
            email: '',
            password: '',
        },
    })

    async function onSubmit(data: SignInForm) {
        setSubmitError(null)

        try {
            await authService.login(data)
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
                    <h1 className="text-lg font-semibold text-slate-100 mb-0.5">Sign in to your account</h1>
                    <p className="text-xs text-[#3a5070] mb-6">Monitor your services 24/7 with real-time alerts</p>
                    <form className="flex flex-col gap-4" onSubmit={handleSubmit(onSubmit)} noValidate>
                        {submitError && (
                            <p className="text-xs text-red-400 bg-red-500/10 border border-red-500/20 rounded-lg px-3 py-2">
                                {submitError}
                            </p>
                        )}
                        <Input
                            id="email"
                            label="Email"
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
                            autoComplete="current-password"
                            error={errors.password?.message}
                            {...register('password')}
                        />
                        <Button type="submit" disabled={isSubmitting}>
                            {isSubmitting ? 'Signing in…' : 'Sign in'}
                        </Button>
                    </form>
                    <p className="text-xs text-[#3a5070] text-center mt-5">Don't have an account? <LinkButton to="/signup">Register</LinkButton></p>
                </div>
                <p className="text-center text-xs text-[#2a3e58] mt-6">
                    Self-hosted · Open source · No telemetry
                </p>
            </div>
        </div>
    )
}
