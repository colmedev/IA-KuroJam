import './CtA.css'
import { Button } from "../Button/Button";

export const Cta = () => {

  return (
    <>
    <section className="cta__container">
      <h2 className='cta__title'>Start your journey</h2>
      <h3 className='cta__subtitle'>Take the first step towards your dream career with our AI-powered guidance.</h3>

      <Button
      text="Get Started"
      link="/app"
      />

    </section>
    </>
  )
}
