import LinkButton from '../../components/LinkButton'

export default function SignUp() {
    return (
        <div className="flex flex-col items-center justify-center h-screen gap-4">
            <h1 className="text-lg font-semibold text-slate-100">Sign Up</h1>
            <p className="text-xs text-[#3a5070]">
                Already have an account? <LinkButton to="/signin">Sign in</LinkButton>
            </p>
        </div>
    )
}
