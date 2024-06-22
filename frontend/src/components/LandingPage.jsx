import React from "react";

const LandingPage = () => {
  //   const history = useHistory();

  return (
    <section className="relative bg-spotifyblack text-white h-screen flex items-center overflow-hidden">
      <img
        src="bg.jpg"
        alt="Background"
        className="absolute  top-0 left-0 w-full h-full  object-cover opacity-10 z-0"
      />

      <div className="mx-auto z-10 max-w-screen-xl px-4 py-32 lg:flex lg:h-screen lg:items-center">
        <div className="mx-auto max-w-3xl text-center">
          <h1 className="bg-gradient-to-r from-spotifygreen via-spotifygreen to-spotifygreen bg-clip-text text-3xl font-extrabold text-transparent sm:text-5xl">
            Welcome to SpotiMood!
          </h1>

          <p className="mx-auto mt-4 max-w-xl sm:text-xl/relaxed font-semibold">
            Happy or sad or relaxed or romantic?
            <span className="sm:block">
              Let SpotiMood help you find the right playlist for your mood.
            </span>
          </p>

          <div className="mt-8 flex flex-wrap justify-center gap-4">
            <a
              onClick={() =>
                (window.location.href = "http://localhost:8080/login")
              }
              className="block w-full rounded border font-bold border-spotifygreen bg-spotifygreen px-12 py-3 text-sm  text-black hover:bg-transparent hover:text-white focus:outline-none focus:ring active:text-opacity-75 sm:w-auto animate-bounce"
              href="#"
            >
              Get Started
            </a>
          </div>
        </div>
      </div>
    </section>
  );
};

export default LandingPage;
