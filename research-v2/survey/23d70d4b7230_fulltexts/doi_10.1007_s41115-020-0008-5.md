## REVIEW ARTICLE

## Hydrodynamics of core-collapse supernovae and their progenitors

## Bernhard Mu ¨ller 1

<!-- image -->

Received: 26 November 2019 / Accepted: 18 April 2020 / Published online: 8 June 2020 /C211 The Author(s) 2020

## Abstract

Multi-dimensional fluid flow plays a paramount role in the explosions of massive stars as core-collapse supernovae. In recent years, three-dimensional (3D) simulations of these phenomena have matured significantly. Considerable progress has been made towards identifying the ingredients for shock revival by the neutrino-driven mechanism, and successful explosions have already been obtained in a number of selfconsistent 3D models. These advances also bring new challenges, however. Prompted by a need for increased physical realism and meaningful model validation, supernova theory is now moving towards a more integrated view that connects multi-dimensional phenomena in the late convective burning stages prior to collapse, the explosion engine, and mixing instabilities in the supernova envelope. Here we review our current understanding of multi-D fluid flow in core-collapse supernovae and their progenitors. We start by outlining specific challenges faced by hydrodynamic simulations of core-collapse supernovae and of the late convective burning stages. We then discuss recent advances and open questions in theory and simulations.

Keywords Supernovae /C1 Massive stars /C1 Hydrodynamics /C1 Convection /C1 Instabilities /C1 Numerical methods

## Contents

|   1 | Introduction.............................................................................................................................   |   3 |
|-----|---------------------------------------------------------------------------------------------------------------------------------------------|-----|
| 1.1 | The multi-dimensional nature of the explosion mechanism........................................                                             |   3 |
| 1.2 | The multi-D structure of supernova progenitors ..........................................................                                   |   5 |
| 1.3 | Observational evidence for multi-D effects in core-collapse supernovae...................                                                   |   6 |
| 1.4 | Scope and structure of this review ...............................................................................                          |   6 |

&amp; Bernhard Mu ¨ller bernhard.mueller@monash.edu

1 School of Physics and Astronomy, Monash University, Clayton, Australia

(0123456789().,-volV

(0123456789().,-volV

<!-- image -->

<!-- image -->

| 2 Numerical methods.................................................................................................................   | 2 Numerical methods.................................................................................................................   | 2 Numerical methods.................................................................................................................   |   7 |
|----------------------------------------------------------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------|-----|
| 2.1 Hydrodynamics ..............................................................................................................       | 2.1 Hydrodynamics ..............................................................................................................       | 2.1 Hydrodynamics ..............................................................................................................       |   8 |
|                                                                                                                                        | 2.1.1                                                                                                                                  | Problem geometry and choice of grids ........................................................                                          |   9 |
|                                                                                                                                        | 2.1.2                                                                                                                                  | Challenges of subsonic turbulent flow.........................................................                                         |  19 |
|                                                                                                                                        | 2.1.3                                                                                                                                  | High-Mach number flow...............................................................................                                   |  23 |
| 2.2                                                                                                                                    | Treatment of gravity......................................................................................................             | Treatment of gravity......................................................................................................             |  25 |
|                                                                                                                                        | 2.2.1                                                                                                                                  | Hydrostatic balance and conservation properties.........................................                                               |  25 |
|                                                                                                                                        | 2.2.2                                                                                                                                  | Treatment of general relativity .....................................................................                                  |  26 |
|                                                                                                                                        | 2.2.3                                                                                                                                  | Poisson solvers ..............................................................................................                         |  27 |
| 2.3                                                                                                                                    | Reactive flow.................................................................................................................         | Reactive flow.................................................................................................................         |  28 |
| 3 Late-stage convective burning in supernova progenitors .....................................................                         | 3 Late-stage convective burning in supernova progenitors .....................................................                         | 3 Late-stage convective burning in supernova progenitors .....................................................                         |  31 |
| 3.1                                                                                                                                    | Interior flow...................................................................................................................       | Interior flow...................................................................................................................       |  32 |
| 3.2                                                                                                                                    | Supernova progenitor models........................................................................................                    | Supernova progenitor models........................................................................................                    |  38 |
| 3.3                                                                                                                                    | Convective boundaries...................................................................................................               | Convective boundaries...................................................................................................               |  39 |
| 3.4                                                                                                                                    | Current and future issues...............................................................................................               | Current and future issues...............................................................................................               |  44 |
| 4 Core collapse and shock revival............................................................................................          | 4 Core collapse and shock revival............................................................................................          | 4 Core collapse and shock revival............................................................................................          |  45 |
| 4.1                                                                                                                                    | Structure of the accretion flow and runaway conditions in spherical symmetry .......                                                   | Structure of the accretion flow and runaway conditions in spherical symmetry .......                                                   |  49 |
| 4.2                                                                                                                                    | Impact of multi-dimensional effects on shock revival.................................................                                  | Impact of multi-dimensional effects on shock revival.................................................                                  |  52 |
| 4.3                                                                                                                                    | Neutrino-driven convection in the gain region ............................................................                             | Neutrino-driven convection in the gain region ............................................................                             |  54 |
| 4.4                                                                                                                                    | The standing accretion shock instability ......................................................................                        | The standing accretion shock instability ......................................................................                        |  58 |
| 4.5                                                                                                                                    | Perturbation-aided explosions .......................................................................................                  | Perturbation-aided explosions .......................................................................................                  |  61 |
| 4.6                                                                                                                                    | Outlook: rotation and magnetic fields in neutrino-driven explosions.........................                                           | Outlook: rotation and magnetic fields in neutrino-driven explosions.........................                                           |  64 |
| 4.7                                                                                                                                    | Proto-neutron star convection and LESA instability....................................................                                 | Proto-neutron star convection and LESA instability....................................................                                 |  66 |
| 5                                                                                                                                      | The explosion phase...............................................................................................................     | The explosion phase...............................................................................................................     |  70 |
| 5.1                                                                                                                                    | The early explosion phase.............................................................................................                 | The early explosion phase.............................................................................................                 |  70 |
| 5.2                                                                                                                                    | Explosion energetics......................................................................................................             | Explosion energetics......................................................................................................             |  72 |
| 5.3                                                                                                                                    | Compact remnant properties .........................................................................................                   | Compact remnant properties .........................................................................................                   |  74 |
| 5.4                                                                                                                                    | Mixing instabilities in the envelope..............................................................................                     | Mixing instabilities in the envelope..............................................................................                     |  77 |

6

References...............................................................................................................................

## 1 Introduction

The death of massive stars is invariably spectacular. In the cores of these stars, nuclear fusion proceeds all the way to the iron group through a sequence of burning stages. At the end of the star's life, nuclear energy generation has ceased in the degenerate Fe core (i.e., the core has become ''inert''), while nuclear burning continues in shells composed of progressively light nuclei further out. Once the core approaches its effective Chandrasekhar mass and becomes sufficiently compact, electron captures on heavy nuclei and partial nuclear disintegration lead to collapse on a free-fall time scale, leaving behind a neutron star or black hole. In most cases, a small fraction of the potential energy liberated during collapse is transferred to the stellar envelope, which is expelled in a powerful explosion known as a corecollapse supernova , as first recognized by Baade and Zwicky (1934).

How precisely the envelope is ejected has remained one of the foremost questions in computational astrophysics ever since the first modeling attempts in the 1960s (Colgate et al. 1961; Colgate and White 1966). In this review we focus on the critical role of multi-dimensional (multi-D) fluid flow during the supernova explosion itself and the final pre-collapse stages of their progenitors.

<!-- image -->

82

For pedagogical reasons, it is preferable to commence our brief exposition of multi-D hydrodynamic effects with the supernova explosion mechanism rather than to follow the sequence of events in nature, or historical chronology.

## 1.1 The multi-dimensional nature of the explosion mechanism

In principle, the collapse of an iron core to a neutron star opens a reservoir of several 10 53 erg of potential energy, which appears more than sufficient to account for the typical inferred kinetic energies of observed core-collapse supernovae of order 10 51 erg (see, e.g., Kasen and Woosley 2009; Pejcha and Prieto 2015).

Transferring the requisite amount of energy from the young ''proto-neutron star'' (PNS) is not trivial, however. The simplest idea is that the energy is delivered when the collapsing core overshoots nuclear density and ' 'bounces'' due to the high incompressibility of matter above nuclear saturation density, which launches a shock wave into the surrounding shells (Colgate et al. 1961). However, the shock wave stalls within milliseconds as nuclear dissociation of the shocked material and neutrino losses drains its initial kinetic energy (e.g., Mazurek 1982; Burrows and Lattimer 1985; Bethe 1990). The shock then turns into an accretion shock, whose radius is essentially determined by the pre-shock ram pressure and the condition of hydrostatic equilibrium between the shock and the PNS surface. It typically reaches a radius between 100 and 200 km a few tens of milliseconds after bounce and then recedes again.

Among the various ideas to ''revive'' the shock (for a more exhaustive overview see Mezzacappa 2005; Kotake et al. 2006; Janka 2012; Burrows 2013) the neutrinodriven mechanism is the most promising scenario and has been explored most comprehensively since it was originally conceived-in a form rather different from the modern paradigm-by Colgate and White (1966). The modern version of this mechanism is illustrated in Fig. 1b: a fraction of the neutrinos emitted from the PNS and the cooling layer at its surface are reabsorbed further out in the ''gain region''. If the neutrino heating is sufficiently strong, the increased thermal pressure drives the shock out, and the heating powers an outflow of matter in its wake.

However, according to the most sophisticated spherically symmetric (1D) models using Boltzmann neutrino transport to accurately model neutrino heating and cooling, this mechanism does not work in 1D (Liebendo ¨rfer et al. 2001; Rampp and Janka 2000) except for the least massive 1 supernova progenitors. For all other progenitors, it is crucial that multi-D effects support the neutrino heating. Convection occurs in the gain region because neutrino heating establishes a negative entropy gradient (Bethe 1990), and was shown to be highly beneficial for obtaining neutrino-driven explosions by the first generation of multi-D models from the 1990s (Herant et al. 1992, 1994; Burrows et al. 1995; Janka and Mu ¨ller 1995, 1996). Another instability, the standing-accretion shock instability (SASI; Blondin et al. 2003; Blondin and Mezzacappa 2006; Foglizzo et al. 2006, 2007;

1 More precisely, the neutrino-driven mechanism works in 1D for progenitors with a steeply declining density profile outside the core (Mu ¨ller 2016) as in the case of electron-capture supernovae from superAGB stars (Kitaura et al. 2006) or the least massive iron core progenitors (Melson et al. 2015b).

<!-- image -->

Fig. 1 Overview of the multi-D effects operating prior to and during a core-collapse supernova as discussed in this review (below the dashed line) in the broader context of the evolution of a massive star. a After millions of years of H and He burning, the star enters the neutrino-cooled burning stages (C, Ne, O, Si burning). These advanced core and shell burning stages tpyically proceed convectively because burning and neutrino cooling at the top and bottom of an active shell/core establish an unstable negative entropy gradient. The interaction of the flow with convective boundaries can lead to mixing and transfer energy and angular momentum by wave excitation. Rotation may modify the flow dynamics. b After the star has formed a sufficiently massive iron core, the core undergoes gravitational collapse, and a young proto-neutron star is formed. The shock wave launched by the ''bounce'' of the core quickly stalls, and is likely revived by neutrino heating in most cases. In the phase leading up to shock revival, neutrino heating drives convection in the heating or ''gain'' region, and the shock may execute large-scale oscillations due the standing accretion shock instability (SASI). Rotation and the asymmetric infall of convective burning shells can modify the dynamics. There is also a convective region below the cooling layer at the protoneutron star (PNS) surface. c After the shock has been revived and sufficient energy has been pumped into the explosion by neutrino heating or some other mechanism, the shock propagates through the outer shells on a time scale of hours to days. During this phase, the interaction of the (deformed) shock with shell interfaces as well as reverse shock formation trigger mixing by the Rayleigh-Taylor (RT), the Richtmyer-Meshkov (RM) instability, and (as a secondary process) the Kelvin-Helmholtz (KH) instability. Once the shock breaks out through the stellar surface, the explosion becomes visible as an electromagnetic transient. Mixing instabilities continue to operate on longer time scales throughout the evolution of the supernova remnant

<!-- image -->

Laming 2007), arises due to an advective-acoustic amplification cycle and gives rise to large-scale ( ' ¼ 1 ; 2) shock oscillations; it plays a similarly beneficial role in the neutrino-driven mechanism as convection (Scheck et al. 2008; Mu ¨ller et al. 2012a). Rapid rotation could also modify the dynamics in the supernova core and support the development of neutrino-driven explosions (Janka et al. 2016; Summa et al. 2018; Takiwaki et al. 2016).

Evidently, multi-D effects are also at the heart of the most serious alternative to neutrino-driven explosions, the magnetohydrodynamic (MHD) mechanism (e.g., Akiyama et al. 2003; Burrows et al. 2007a; Winteler et al. 2012; Mo ¨sta et al. 2014b), which likely explains unusually energetic ''hypernovae''. But whether corecollapse supernovae are driven by neutrinos or magnetic fields, it is pertinent to ask:

<!-- image -->

How important are the initial conditions for the multi-D flow dynamics that leads to shock revival?

## 1.2 The multi-D structure of supernova progenitors

For pragmatic reasons, supernova models have long relied on 1D stellar evolution models as input, or at best on ' '1.5D'' rotating models using the shellular approximation (Zahn 1992). For non-rotating progenitors, spherical symmetry was either broken by introducing perturbations in supernova simulations by hand, or due to grid perturbations. For rotating and magnetized progenitors, spherical symmetry is broken naturally, but on the other hand the stellar evolution models do not provide the detailed multi-D angular momentum distribution and magnetic field geometry, which must be specified by hand.

In reality, even non-rotating progenitors are not spherically symmetric at the onset of collapse. Outside the iron core, there are typically several active convective burning shells (Fig. 1a) that will collapse in the wake of the iron core within hundreds of milliseconds. It was realized in recent years that the infall of asymmetric shells can be important for shock revival (Couch and Ott 2013; Mu ¨ller 2015; Couch et al. 2015; Mu ¨ller et al. 2017a).

The multi-D structure of supernova progenitors is thus directly relevant for the neutrino-driven mechanism, but the potential ramifications of multi-D effects during the pre-collapse phase are in fact much broader: How do they affect mixing at convective boundaries, and hence the evolution of the shell structure on secular time scales? How do they affect the angular momentum distribution and magnetic fields in supernova progenitors?

## 1.3 Observational evidence for multi-D effects in core-collapse supernovae

Observations contain abundant clues about the multi-D nature of core-collapse supernova explosion. Large birth kicks of neutron stars (Hobbs et al. 2005; FaucherGigue `re and Kaspi 2006; Ng and Romani 2007) and even black holes (Repetto et al. 2012) cannot be explained by stellar dynamics alone and require asymmetries in the supernova engine. There is also evidence for mixing processes during the explosion and large-scale asymmetries in the ejecta from the spectra and polarization signatures of many observed transients (e.g., Wang and Wheeler 2008; Patat 2017), and from young supernova remnants like Cas A (Grefenstette et al. 2014).

The relation between the asymmetries in the progenitor and the supernova core, and the asymmetries in observed transients and gaseous remnants is not straightforward, however. The observable symmetries are rather shaped by mixing processes that operate as the shock propagates through the stellar envelope (Fig. 1c). Rayleigh-Taylor instability occurs behind the shock as it scoops up material and decelerates (Chevalier 1976; Bandiera 1984), and the interaction of a non-spherical shock with shell interfaces can give rise to the Richtmyer-Meshkov instability (Kifonidis et al. 2006). The asymmetries imprinted during the first seconds of an explosion provide the seed for these late-time mixing instabilities, and 3D supernova modellling is now moving towards an integrated approach from the

<!-- image -->

early to the late stages of the to better link the observations to the physics of the explosion mechanism (e.g., Hammer et al. 2010; Wongwathanarat et al. 2013; Mu ¨ller et al. 2018; Chan et al. 2018, 2020), and, in future, even to the multi-D progenitor structure.

## 1.4 Scope and structure of this review

In this review, we are primarily concerned with the numerical techniques for modeling multi-D fluid flow in core-collapse supernovae and their progenitors, and with our current understanding of the theory and phenomenology of the multi-D fluid instabilities. Although multi-D effects are relevant to virtually all aspects of core-collapse supernova theory, we can only afford cursory attention to many of them in order to keep this overview focused.

There are many other reviews to fill the gap, or provide a different perspective. Janka (2012) provides a very broad, but less technical, overview of the core-collapse supernova problem. For the explosion mechanism and a different take on the role of multi-D effects, the reader may also consult Mezzacappa (2005), Kotake et al. (2006), Burrows (2013), Mu ¨ller (2016), Janka et al. (2016) and Couch (2017). The important problem of neutrino transport is treated in considerable depth by Mezzacappa (2005, 2020) and is therefore not addressed in this review. A number of reviews address the potential of neutrinos (e.g., Kotake et al. 2006; Mu ¨ller 2019b) and gravitational waves (Ott 2009; Kotake 2013; Kalogera et al. 2019) as diagnostics of the multi-dimensional dynamics in the supernova core.

We shall start by discussing the governing equations for reactive, self-gravitating flow and their numerical solution in the context of core-collapse supernovae and convective burning in Sect. 2. We do not treat numerical methods for MHD in supernova simulations, although we occasionally comment on the role of MHD effects in the later sections. In the subsequent sections, we then review recent simulation results and progress in the theoretical understanding of convection during the late burning stages (Sect. 3), of supernova shock revival (Sect. 4), and the hydrodynamics of the explosion phase (Sect. 5).

## 2 Numerical methods

Modeling the late stages of nuclear burning and the subsequent supernova explosion involves solving the familiar equations for mass, momentum, and energy conservation with source terms that account for nuclear burning and the exchange of energy and momentum with neutrinos. Viscosity and thermal heat conduction mediated by photons, electron/positrons, and ions can be neglected, and so we have (in the Newtonian limit and neglecting magnetic fields),

<!-- formula-not-decoded -->

<!-- image -->

<!-- formula-not-decoded -->

<!-- formula-not-decoded -->

in terms of the density q , the fluid velocity v , the pressure P , internal energy density e , the gravitational potential U , and the neutrino energy and momentum source terms Q e and Q m. If we take e to include nuclear rest-mass contributions, there is no source term for the nuclear energy generation rate; otherwise an additional source term \_ Q nuc appears on the right-hand side (RHS) of Eq. (3). These equations are supplemented by conservation equations for the mass fractions Xi of different nuclear species and the electron fraction Y e (net number of electrons per baryon),

<!-- formula-not-decoded -->

<!-- formula-not-decoded -->

where the source terms \_ Xi ; burn and \_ QY e account for nuclear reactions and the change of the electron fraction by b processes.

In the regime of sufficiently high optical depth, the effect of the neutrino source terms could alternatively be expressed by non-ideal terms for heat conduction, viscosity, and diffusion of lepton number (e.g., Imshennik and Nadezhin 1972; Bludman and van Riper 1978; Goodwin and Pethick 1982; van den Horn and van Weert 1983, 1984; Yudin and Nadyozhin 2008), but this approach would break down at low optical depth. The customary approach to Eqs. (1-5) is, therefore, to apply an operator-split approach and combine a solver for ideal hydrodynamics for the left-hand side (LHS) and the gravitational source terms with separate solvers for the source terms due to neutrino interactions and nuclear reactions. Simulations of the Kelvin-Helmholtz cooling phase of the PNS over time scales of seconds form an exception; here only the PNS interior is of interest so that it is possible and useful to formulate the neutrino source terms in the equilibrium diffusion approximation (Keil et al. 1996; Pons et al. 1999).

## 2.1 Hydrodynamics

A variety of computational methods are employed to solve the equations of ideal hydrodynamics in the context of supernova explosions or the late stellar burning stages. Nowadays, the vast majority of codes use Godunov-based high-resolution shock capturing (HRSC) schemes with higher-order reconstruction (see, e.g., LeVeque 1998b; Martı ´ and Mu ¨ller 2015; Balsara 2017 for a thorough introduction). Examples include implementations of the piecewise parabolic method of Colella and Glaz (1985) or extensions thereof in the Newtonian hydroydnamics codes PROMETHEUS (Fryxell et al. 1989, 1991; Mu ¨ller et al. 1991), which has been integrated into various neutrino transport solvers by the Garching group (Rampp

<!-- image -->

3

and Janka 2002; Buras et al. 2006b; Scheck et al. 2006), its offshoot PROMPI (Meakin and Arnett 2007b), PPMSTAR (Woodward et al. 2019), VH1 (Blondin and Lufkin 1993; Hawley et al. 2012) as used within the CHIMERA transport code (Bruenn et al. 2018), FLASH (Fryxell et al. 2000), CASTRO (Almgren et al. 2006), ALCAR (Just et al. 2015), and FORNAX (Skinner et al. 2019). This approach is also used in most general relativistic (GR) hydrodynamics codes for core-collapse supernovae like COCONUT (Dimmelmeier et al. 2002; Mu ¨ller et al. 2010), ZELMANI (Reisswig et al. 2013; Roberts et al. 2016), and GRHYDRO (Mo ¨sta et al. 2014a). Godunov-types scheme with piecewise-linear total-variation diminishing (TVD) reconstruction are still used in the FISH code (Ka ¨ppeli et al. 2011) and in the relativistic FUGRA code of (Kuroda et al. 2012, 2016b). The 3DnSNe code of the Fukuoka group (e.g., Takiwaki et al. 2012), which is based on the ZEUS code of Stone and Norman (1992), has also switched from an artificial viscosity scheme to a Godunov-based finite-volume approach with TVD reconstruction (Yoshida et al. 2019).

Alternative strategies are less frequently employed. The VULCAN code (Livne 1993; Livne et al. 2004) uses a staggered grid and von Neumann-Richtmyer artificial viscosity (Von Neumann and Richtmyer 1950). The SNSPH code (Fryer et al. 2006) uses smoothed-particle hydrodynamics (Gingold and Monaghan 1977; Lucy 1977; for modern reviews see Price 2012; Rosswog 2015). Although less widely used in supernova modeling today, the SPH approach has been utilized for some of the early studies of Rayleigh-Taylor mixing (Herant and Benz 1991) and convectively-driven explosions (Herant et al. 1992) in 2D and later for the first 3D supernova simulations with gray neutrino transport (Fryer and Warren 2002). Multidimensional moving mesh schemes have been occasionally employed to simulate magnetorotational supernovae (Ardeljan et al. 2005) and reactive-convective flow in stellar interiors (Dearborn et al. 2006; Stancliffe et al. 2011). More recently ' 'second-generation'' moving-mesh codes based on Voronoi tessellation, such as TESS (Duffell and MacFadyen 2011) and AREPO (Springel 2010), have been developed and employed for simulations of jet outflows (Duffell and MacFadyen 2013, 2015), and fallback supernovae (Chan et al. 2018). Spectral solvers for the equations of hydrodynamics, while popular for solar convection, have so far been applied only once for simulating oxygen burning (Kuhlen et al. 2003), but never for the core-collapse supernova problem.

Since Godunov-based finite-volume solvers are now most commonly employed for simulating core-collapse supernovae and the late burning stages, we shall focus on the problem-specific challenges for this approach in this section.

## 2.1.1 Problem geometry and choice of grids

The physical problem geometry in global simulations of core-collapse supernovae and the late convective burning stages is characterized by approximate spherical symmetry, and one frequently needs to deal with strong radial stratification and a large range of radial scales. For example, during the pre-explosion and early explosion phase, the PNS develops a ' 'density cliff'' at its surface that is

<!-- image -->

approximately in radiative equilibrium and can be approximated as an exponential isothermal atmosphere with a scale height H of

<!-- formula-not-decoded -->

in terms of the PNS mass M , radius R , and surface temperature T , and the baryon mass m b. With typical values of M /C24 1 : 5 M /C12 , R shrinking down to a final value of /C24 12 km, and a temperature of a few MeV, the scale height soon shrinks to a few 100 m. Later on during the explosion, the scales of interest shift to the radius of the entire star, which is of order /C24 10 8 km for red supergiants.

The spherical problem geometry and the multi-scale nature of the flow is a critical element in the choice of the numerical grid for ''star-in-a-box'' simulations. 2 Cartesian grids, various spherical grids, and, on occasion, unstructured grids have been used in the field for global simulations and face different challenges.

Grid-induced perturbations Cartesian grids have the virtue of algorithmic simplicity and do not suffer from coordinate singularities, but also come with disadvantages as they are not adapted to the approximate symmetry of the physical problem. The unavoidable non-spherical perturbations from the grid make it impossible to reproduce the spherically symmetric limit in multi-D even for perfectly spherical initial conditions, or to study the growth of non-spherical perturbations in a fully controlled manner. The former deficiency is arguably an acceptable sacrifice, though it can limit the possibilities for code testing and verification, but the latter can introduce visible artifacts in simulations. For example, Cartesian codes sometimes produce non-vanishing gravitational wave signals from the bounce of non-rotating cores (Scheidegger et al. 2010), and often show dominant ' ¼ 4 modes during the growth phase of non-radial instabilities (Ott et al. 2012).

Handling the multi-scale problem Furthermore, a single Cartesian grid cannot easily handle the multiple scales encountered in the supernova problem. Even with Oð 1000 3 Þ zones, such a grid can at best cover the region inside 1000 km with acceptable resolution, but following the infall of matter for several 100 ms without boundary artifacts and the development of an explosion requires covering a region of at least 10 ; 000 km. This problem is often dealt with by using adaptive mesh refinement (AMR; see, e.g., Berger and Colella 1989; Fryxell et al. 2000), which is usually implemented as ' 'fixed mesh refinment'' for pre-defined nested cubic patches (e.g., Schnetter et al. 2004). Other codes have opted to combine a single central Cartesian patch or nested patches with a spherically symmetric region (Scheidegger et al. 2010) or multiple spherical polar patches (Ott et al. 2012) outside. For long-time simulations of Rayleigh-Taylor mixing in the envelope, standard adaptive or pre-defined mesh refinement may not be sufficiently efficient for covering the range of changing scales and necessitate manual remapping to a

2 Of course, some problems can or need to be studied using simplified geometries (planar or cylindrical) or local simulations.

<!-- image -->

coarser grids (''homographic expansion'') as the simulation proceeds (Chen et al. 2013).

In spherical coordinates, the multi-scale nature of the problem can be accommodated to a large degree by employing a non-uniform radial grid that transitions to roughly equal spacing in log r at large radii. Radial resolution can be added selectively in strongly stratified regions like the PNS surface (Buras et al. 2006b), or one can use an adaptive moving radial grid (Liebendo ¨rfer et al. 2004; Bruenn et al. 2018). However, some care must be exercised in using non-uniform radial grids. Rapid variations in the radial grid resolution D r = r can produce artifacts such as artificial waves and disturbances of hydrostatic equilibrium.

It is also straightforward to implement a moving radial grid to adapt to changing resolution requirements or the bulk contraction/expansion of the region of interest; see Winkler et al. (1984), Mu ¨ller (1994) for an explanation of this technique. The MPA/Monash group routinely apply such a moving radial grid in quasi-Lagrangian mode during the collapse phase (Rampp and Janka 2000), and, with a prescribed grid function, in parameterized multi-D simulations of neutrino-driven explosions (Janka and Mu ¨ller 1996; Scheck et al. 2006) and in simulations of convective burning (Mu ¨ller et al. 2016b). The Oak Ridge group uses a truly adaptive radial grid in their supernova simulations with the CHIMERA code (Bruenn et al. 2018). A moving radial mesh might also appear useful for following the expansion of the ejecta and the formation of a strongly diluted central region in simulations of mixing instabilities in the envelope, but the definition of an appropriate grid function is nontrivial. Most simulations of mixing instabilities in spherical polar coordinates have therefore relied on simply removing zones continuously from the evacuated region of the blast wave to increase the time step (Hammer et al. 2010) rather than implementing a moving radial mesh (Mu ¨ller et al. 2018).

Both fixed mesh refinement and spherical grids with non-uniform radial mesh spacing only provide limited adaptability to the structure of the flow. Truly adaptive mesh refinement can provide superior resolution in cases where very tenuous, non volume-filling flow structures emerge. Mixing instabilities in the envelope are a prime example for such a situation, and have often been studied using AMR in spherical polar coordinates (Kifonidis et al. 2000, 2003, 2006) in 2D and Cartesian coordinates (Chen et al. 2017).

The time step constraint in spherical polar coordinates While spherical polar coordinates are well-adapted to the problem geometry, they also suffer from drawbacks. One of these drawbacks-among others that we discuss further belowis that the converging cell geometry imposes stringent constraints on the time step near the grid axis and the origin. The Courant-Friedrichs-Lewy limit (Courant et al. 1928) for the time step D t requires that D t \ r Dh = ðj v j þ c s Þ in 2D and D t \ r sin h Du = ðj v j þ c s Þ in 3D in terms of the grid spacing Dh and Du in latitude and longitude, and the fluid velocity v and sound speed c s . If Dh ¼ Du , this is worse than a Cartesian code with grid spacing comparable to D r by a factor Dh /C28 1 in 2D and Dh 2 /C28 1 in 3D near the origin.

Various workarounds have been developed to tame this time-step constraint. Some core-collapse supernova codes (PROMETHEUS-VERTEX, CHIMERA, COCONUT)

<!-- image -->

simulate the innermost region of the grid assuming spherical symmetry. The approximation of spherical symmetry is well justified in the core, since the innermost region of the PNS is convectively stable during the first seconds after collapse and explosion until the late Kelvin-Helmholtz cooling phase. Even more savings can be achieved by treating the PNS convection zone using mixing-length theory (Mu ¨ller 2015), but this is a more severe approximation that significantly affects the predicted gravitational wave signals and certain features of the neutrino emission and nucleosynthesis. Concerns have also been voiced that the imposition of a spherical core region creates an immobile obstacle to the flow that leads to the violation of momentum conservation, which might have repercussions on neutron star kicks (Nordhaus et al. 2010a). While it is true that the PNS tends not to move in simulations with a spherical core region, Scheck et al. (2006) found (using a careful analysis based on hydro simulations in the accelerated frame comoving with the PNS) that the assumption of an immobile PNS does not gravely affect the dynamics in the supernova interior and the PNS kick in particular.

Even with a 1D treatment for the innermost grid zones, one is still left with a severe time-step constraint at the grid axis in 3D. A number of alternatives to spherical polar grids with uniform spacing in latitude h and longitude u can help to remedy this. The simplest workaround is to adopt uniform spacing in l ¼ cos h instead of h . In this case, one has sin h ¼ ð 2 N h /C0 1 Þ 1 = 2 = N h /C25 ffiffi ffi 2 p N /C0 1 = 2 h in the zones adjacent to the axis for N h zones in latitude instead of sin h /C25 N /C0 1 h = 2, so the time step limit scales as D t / ffiffi ffi 2 p N /C0 1 = 2 h N /C0 1 / instead of D t / N /C0 1 h N /C0 1 / = 2, where N u is the number of zones in longitude. Alternatively, one can selectively increase the h -grid spacing in the zones close to the axis. However, the time step constraint at the axis is still more restrictive than at the equator in this approach, and the aspect ratio of the grid cells becomes extreme near the pole, which can create problems with numerical stability and accuracy.

One approach to fully cure the time step problem, which was first proposed for simulations of compact objects by Cerda ´-Dura ´n (2009), consists in abandoning the logically Cartesian grid in r , h , and u and selectively coarsening the grid spacing in u (and possibly h ) near the axis (and optionally at small r ) as illustrated in Fig. 2. Such a mesh coarsening scheme has been included in the COCONUT-FMT code (Mu ¨ller 2015) with coarsening in the u -direction, and as a ''dendritic grid'' with coarsening in the h - and u -direction in the FORNAX code (Skinner et al. 2019) and the 3DnSNe code (Nakamura et al. 2019). Mesh coarsening can be implemented following standard AMR practice by prolongating from the coarser grids to the finer grids in the reconstruction step. Alternatively, one can continue using the hydro solver on a fine uniform grid in h and u , and average the solution over coarse ' 'supercells'' after each time step, followed by a conservative prolongation or ''prereconstruction'' step back onto the fine grid to ensure higher-order convergence. This has the advantage of retaining the data layout and algorithmic structure of a spherical polar code, but care is required to to ensure that the prolongation of the conserved variables does not introduce non-monotonicities in the primitive variables, which limits the pre-reconstruction step to second-order accuracy in practice (Mu ¨ller et al. 2019). A possible concern with mesh coarsening on standard

<!-- image -->

Fig. 2 Alternative spherical grids that avoid the tight time step constraint at the axis of standard spherical polar grids: a Grid with mesh coarsening in the u -direction only. Only an octant of the entire grid is shown. b Dendritic grid with coarsening in the h - and u -direction Image reproduced with permission from (Skinner et al. 2019), copyright by AAS. c Overset Yin-Yang grid (Kageyama and Sato 2004; Wongwathanarat et al. 2010a) with two overlapping spherical polar patches in yellow and cyan

<!-- image -->

spherical polar grids is that it may favor the emergence of axis-aligned bipolar flow structures during the explosion phase in supernova simulations (Mu ¨ller 2015; Nakamura et al. 2019). In practice, however, strong physical seed perturbations easily break any grid-induced alignment of the flow with the axis. As the more simulations with mesh coarsening become available (Mu ¨ller et al. 2019; Burrows et al. 2020), it does not appear that axis alignment is a recurring problem.

Filtering in Fourier space, which has long been used in the meteorology community (Boyd 2001), provides another means of curing the restrictive time step constraint near the axis and has been implemented in the COCONUT-FMT code (Mu ¨ller et al. 2019). This can also be implemented with minimal interventions in a solver for spherical polar grids, and is attractive because the amount of smoothing that is applied to the solution increases more gradually towards the axis than with mesh coarsening schemes. Mu ¨ller et al. (2019) suggests that this eliminates the problem of axis-aligned flows. On the downside, simulations with Fourier filtering may occasionally encounter problems with the Gibbs phenomenon at the shock.

More radical solutions to the axis problem include overset grids and nonorthogonal spherical grids. An overset Yin-Yang grid (Kageyama and Sato 2004) has been implemented in the PROMETHEUS code (Wongwathanarat et al. 2010a; Melson 2013) and used successfully for simulations of supernovae and convective burning. The Yin-Yang grid provides near-uniform resolution in all directions, solves the time step problem, and also eliminates the delicate problem of boundary conditions at the axis of a spherical polar grid. The added algorithmic complexity is limited to interpolation routines that provide boundary conditions; since each patch is part of a spherical polar grid, no modifications of the hydro solver for nonorthogonal grids are required. As a downside, it is more complicated-but possible (Peng et al. 2006)-to implement overset grids in a strictly conservative manner. In future, non-orthogonal grids spherical grids (Ronchi et al. 1996; Calhoun et al. 2008; Wongwathanarat et al. 2016) may provide another solution that avoids the axis problem and ensures conservation in a straightforward manner, but applications

<!-- image -->

are so far limited to other astrophysical problems (Koldoba et al. 2002; Fragile et al. 2009; Shiota et al. 2010).

Boundary conditions The definition of the outer boundary conditions for Cartesian grids can be more delicate and less flexible than for spherical grids. Simulations of supernova shock revival and the late convective burning stages usually do not cover the entire star for efficiency reasons, and sometimes it can even be desirable to excise an inner core region, e.g., the PNS interior in supernova simulations (e.g., Janka and Mu ¨ller 1996; Scheck et al. 2008; Ugliano et al. 2012; Ertl et al. 2016; Sukhbold et al. 2016) or the Fe/Si core in the O shell burning models of Mu ¨ller et al. (2017a). To minimize artifacts near the outer and (for annular domains) the inner boundary, the best strategy is often to impose hydrostatic boundary conditions assuming constant entropy, so that the pressure P , density q , and radial velocity vr in the ghost cells are obtained as

Z

<!-- formula-not-decoded -->

Z

<!-- formula-not-decoded -->

<!-- formula-not-decoded -->

in terms of the radial gravitational acceleration g and the sound speed c s (cf. Zingale et al. 2002 for hydrostatic extrapolation in the plane-parallel case). This can be readily implemented for spherical grid, and the same is true for inflow, outflow, or wall boundary conditions for excised outer shells or an excised core.

In Cartesian coordinates, however, defining boundary conditions as a function of radius is at odds with the usual strategy of enforcing the boundary conditions by populating ghost zones along individual grid lines separately. For pragmatic reasons, one often enforces standard boundary conditions (reflecting/inflow/outflow) on the faces of the cubical domain instead (e.g., Couch et al. 2015), which is viable as long as the domain boundaries are sufficiently distant from the region of interest. Alternatively, one can impose fixed boundary conditions inside the cubical domain, but outside a smaller spherical region of interest (e.g., Woodward et al. 2018). Outflow conditions on an interior boundary, e.g., for fallback onto a compact remnant, can also be implemented relatively easily (Joggerst et al. 2009).

On the other hand, the boundary conditions at the axis and the origin require careful consideration in case of a spherical polar in order to minimize artifacts from the grid singularities. Conventionally, one uses reflecting boundary conditions to populate the ghost zones before performing the reconstruction in the r - and h -direction, i.e., one assumes odd parity for the velocity components vr and v h respectively, and even parity for scalar quantities and the transverse velocity components. This usually ensures that vr and v h do not blow up near the grid singularities, but in some cases stronger measures are required; e.g., one can enforce zero vr or v h in the cell next to the origin/grid axis, or switch to step function

<!-- image -->

reconstruction in the first cell. One may also need to impose odd parity for v u or for better stability, or reconstruct the angular velocity component xu ¼ v u = r instead of v u .

No hard-and-fast rules for such fixes at the axis and the origin can be given, except perhaps that one should also consider treating the geometric source terms in spherical polar coordinates differently (see below) before applying fixes to the boundary conditions that reduce the order of reconstruction, or before manually damping or zeroing velocity components. In fact, the symmetry assumptions behind reflecting boundary conditions (i.e., vr ! 0 at the origin and v h ! 0) are actually too strong. Strictly speaking, one should only impose the condition that the Cartesian velocity components vx , vy , and vz are continuous across the singularity for smooth flow. In principle, this can be accommodated during the reconstruction by populating the ghost zones for r \ 0, h \ 0, and h [ p with values from the corresponding grid lines across the origin or the axis, bearing in mind any flip in direction of the basis vectors e r , e h , and e u across the coordinate singularity. For the reconstruction along the radial grid line with constant h and u , this comes down to defining

/C26

<!-- formula-not-decoded -->

/C26

<!-- formula-not-decoded -->

/C26

<!-- formula-not-decoded -->

and, similarly, for reconstruction in the h -direction along a grid line with constant r and u :

8

<!-- formula-not-decoded -->

8

<!-- formula-not-decoded -->

<!-- image -->

<!-- formula-not-decoded -->

8

<!-- formula-not-decoded -->

This allows for non-zero values of vr and v h at the origin to reflect that matter can flow across the origin and the axis. Such special polar boundary conditions have been implemented for 3D light-bulb 3 simulations of SASI and convection with FLASH (Ferna ´ndez 2015), and are also used in the FORNAX code (Skinner et al. 2019). In practice, however, reflecting boundary conditions do not appear to pose a major obstacle for flows across the axis or the origin if the diverging fictitious force terms are treated appropriately (see below). The reason is that reflecting boundaries merely slightly degrade the accuracy of the first cell interfaces away from the origin and the axis; the fact that velocity components at the coordinate singularity are (incorrectly) forced to zero on the cell interfaces at r ¼ 0, h ¼ 0, and h ¼ p does not matter much because these interfaces have a vanishing surface area, so that the interface fluxes must vanish anyway.

Geometric source terms Another obstacle in spherical polar coordinates is the occurrence of fictitious force terms in the momentum equation. In terms of the density q and the orthonormal components vi and gi of the velocity and gravitational acceleration, the equations read,

<!-- formula-not-decoded -->

<!-- formula-not-decoded -->

<!-- formula-not-decoded -->

where the fictitious force terms are singular at the origin and at the axis. Often straightforward time-explicit discretization is sufficient for these source terms, especially in unsplit codes with Runge-Kutta time integration. In a dimensionally split implementation it can be advantageous to include a characteristic state correction in the Riemann problem due to (some of the) fictitious force terms (Colella and Woodward 1984). Time-centering of the geometric source terms can also lead to minor differences (Ferna ´ndez 2015).

When stability problems or pronounced axis artifacts are encountered, one can adopt a more radical solution and transport the Cartesian momentum density q v ¼

3 Broadly speaking, light-bulb simulations manually fix the neutrino luminosities and spectral properties, or compute them based on simple analytic considerations, and also use simplified neutrino source terms.

<!-- image -->

ð q vx ; q vy ; q vz Þ while still using the components v a in the spherical polar basis as advection velocities, so that the fictitious force terms disappear entirely, ffiffi

<!-- formula-not-decoded -->

where a 2 f r ; h ; u g and c ¼ r 4 sin 2 h is the determinant of the metric. This has been implemented in the COCONUT-FMT code as one of several options for the solution of the momentum equation. Transforming back and forth between the spherical polar basis for the reconstruction and solution of the Riemann problem and the Cartesian components for the update of the conserved quantities might appear cumbersome, but in fact one need not explicitly transform to Cartesian components at all. Instead one only needs to rotate vectorial quantities from the interface to the cell center when updating the momentum components in the spherical polar basis. For example, for uniform grid spacing Dh in the h -direction, the flux difference terms from the h -interfaces j and j þ 1 for updating q vr and q v h in zone j þ 1 = 2 become

/C18

/C18

/C19

<!-- formula-not-decoded -->

/C19

<!-- formula-not-decoded -->

where D A is the interface area and Dh is the grid spacing in the h -direction, which is assumed to be uniform here. The term for q v u is not modified at all. Apart from eliminating the fictitious force terms in favor of flux flux terms, this alternative discretization of the momentum advection and pressure terms also complies with the conservation of total momentum (although discretization of the gravitational source term may still violate momentum conservation).

The singularities at the origin and the pole constitute a more severe problem for relativistic codes using free evolution schemes for the metric, where they can jeopardize the stability of the metric solver. We refer to Baumgarte et al. (2013, 2015) for a robust solution to this problem that has been implemented in their NADA code; they employ a reference metric formulation both for the field equations and the fluid equations that factors out metric terms that become singular and use a partially implicit Runge-Kutta scheme to evolve the problematic terms.

Angular momentum conservation A somewhat related issue concerns the violation of angular momentum conservation in standard finite-volume codes (both with Eulerian grids and moving meshes). This is a concern especially for problems such as convection in rotating stars and magnetorotational explosions where rotation plays a major dynamical role and the evolution of the flow needs to be followed over long time scales. It is also an issue for question such as pulsar spin-up by asymmetric accretion, although a post-processing of the numerical angular

<!-- image -->

momentum flux can help to obtain meaningful results even when there is a substantial violation of angular momentum conservation (Wongwathanarat et al. 2013).

The problem of angular momentum non-conservation can be solved, or at least mitigated, using Discontinuous Galerkin methods (Despre ´s and Labourasse 2015; Mocz et al. 2014; Schaal et al. 2015), which are not currently used in this field however, and can be avoided entirely in SPH (Price 2012). For a given numerical scheme, increasing the resolution is usually the only solution to minimize the conservation error, but in 3D spherical polar coordinates (and in 2D cylindrical coordinates), one can still ensure exact conservation of the angular momentum component Lz along the grid axis by conservatively discretizing the conservation equation for q v u r sin h ,

<!-- formula-not-decoded -->

instead of Eq. (18). Incidentally, this angular-momentum conserving formulation emerges automatically in GR hydrodynamics in spherical polar coordinates if one solves for the covariant momentum density components as in the COCONUT code (Dimmelmeier et al. 2002). However, as a price for exact conservation of Lz , one occasionally encounters very rapid rotational flow around the axis, and enforcing conservation of only one angular component may add to artificial flow anisotropies due to the spherical polar grid geometry. Moreover, this recipe cannot be used for Yin-Yang-type overset spherical grids or for non-orthogonal spherical grids. If angular momentum conservation is a concern, one can, however, resort to a compromise by conservatively discretizing the equations for q v h r and q v u r ,

<!-- formula-not-decoded -->

<!-- formula-not-decoded -->

which eliminates some of the fictitious force terms. This sometimes considerably improves angular momentum conservation for all angular momentum components and works for any spherical grid. Figure 3 illustrates the difference between the standard form of the fictitious force term and the alternative form in Eqs. (23,24) for a simulation of oxygen shell convection in a rapidly rotating gamma-ray burst progenitor.

<!-- image -->

Fig. 3 Profiles of the mass-weighted, spherically-averaged specific angular momentum h j i in 3D simulations of convective oxygen shell burning in a rapidly rotating gamma-ray burst progenitor with an initial helium core mass of 16 M /C12 from Woosley and Heger (2006). The red and blue curves show the angular momentum distribution after a simulation time of 522 s (about 15 convective turnovers) for the standard formulation (red) of the fictitious force terms in Eqs. (17,18) and for the alternative formulation (blue) in Eqs. (23,24). The initial angular momentum profile is shown in black. The alternative formulation reduces the violation of global angular momentum conservation from 20% to 7%. However, for the standard formulation the effects of angular momentum non-conservation are much bigger locally than suggested by the global conservation error. The fictitious force terms proportional to vr lead to a considerable loss of angular momentum near the reflecting boundaries, even though the angular momentum flux through the boundaries is exactly zero, and there is a spurious increase of angular momentum at the bottom of the convective oxygen burning shell outside the mass coordinate m ¼ 2 : 7 M /C12

<!-- image -->

## 2.1.2 Challenges of subsonic turbulent flow

In the final stages of massive stars one encounters a broad range of flow regimes. The convective Mach number, i.e., the ratio of the typical convective velocity to the sound speed, during the advanced convective burning stages is fairly low, ranging from /C24 10 /C0 3 or less (Cristini et al. 2017) to /C24 0 : 1-0.2 in the innermost shells at the onset of collapse (Collins et al. 2018). The flow is highly turbulent with Reynolds numbers of order 10 13 -10 16 . During the supernova, one finds convective Mach numbers from a few 10 /C0 2 in PNS convection to /C24 0 : 4 (Mu ¨ller and Janka 2015; Summa et al. 2016) in the gain region around shock revival, and as the shock propagates through the envelope, the flow becomes extremely supersonic with Mach numbers of up to several 10 2 . The flow in unstable regions is typically highly turbulent. In the gain region, one obtains a nominal Reynolds number of order 10 17 based on the neutron viscosity (Abdikamalov et al. 2015), but non-ideal effects of neutrino viscosity and drag play a role in the environment of the PNS (Burrows and Lattimer 1988; Guilet et al. 2015; Melson et al. 2020). In the outer regions of the PNS convection zone neutrino viscosity keeps the Reynolds number as low as /C24 100 during some phases (Burrows and Lattimer 1988), and in the gain region drag effects are still so large that the flow cannot be assumed to behave like ordinary high-Reynolds number flow (Melson et al. 2020). Later on, as the shock propagates

<!-- image -->

through the envelope and mixing instabilities develop, neutrino drag becomes unimportant, and the flow is again in the regime of very high Reynolds numbers.

Both the vast range in Mach number and the turbulent nature of the flow present challenges for the accuracy and robustness of numerical simulations. While even relatively simple HRSC schemes can deal with Mach numbers of Ma J 1 with aplomb at reasonably high resolution, they can become excessively diffusive at low Mach number because of spurious acoustic waves arising from the discontinuities in the reconstruction (Guillard and Murrone 2004; Miczek et al. 2015). Moreover, the acoustic time step constraint in explicit finite-volume codes becomes excessively restrictive at low Mach number compared to the physical time scales of interest. These problems can be dealt with by using various low-Mach number approximations, e.g., the anelastic approximation as in Glatzmaier (1984), or more general formulations as in the MAESTRO code of Nonaka et al. 2010, or by a time-implicit discretization of the full compressible Euler equations (Viallet et al. 2011; Miczek et al. 2015). However, few studies (Kuhlen et al. 2003; Michel 2019) have used such methods to deal with low-Mach number flow in the late stages of convective burning massive stars as yet; they have been employed more widely to study, e.g., the progenitors of thermonuclear supernovae (Zingale et al. 2011; Nonaka et al. 2012)

Part of the reason is that advanced HRSC schemes remain accurate and competitive down to Mach numbers of 10 /C0 2 and below depending on the reconstruction method and the Riemann solver. Although it is impossible to decide a priori whether a particular choice of methods is adequate for a given physical problem, or what its resolution requirements are, it is useful to be aware of strengths and weaknesses of different schemes in the context of subsonic turbulent flow. Unfortunately, our discussion of these strengths and weaknesses must remain rather qualitative because very few studies in the field have compared the performance of different Riemann solvers and reconstruction schemes in full-scale simulations and not only for idealized test problems.

Riemann solvers For supernova simulations with Godunov-based codes, a variety of Riemann solvers are currently used. Newtonian codes typically opt either for a (nearly) exact solution of the Riemann problem following Colella and Glaz (1985) 4 or for approximate solvers that at least take the full wave structure of the Riemann problem into account such as the HLLC solver (Toro et al. 1994). On the other hand, the majority of relativistic simulations still resort to the HLLE solver (Einfeldt 1988) because of the added complexity of full-wave approximate Riemann solvers for GR hydrodynamics; exceptions include the COCONUT code which routinely uses the relativistic HLLC solver of Mignone and Bodo (2005), the COCONUT simulations of Cerda ´-Dura ´n et al. (2005) using the Marquina solver (Donat and Marquina 1996), and the convection simulations with the WHISKEYTHC code (Radice et al. 2016), which uses a Roe-based flux-split scheme.

4 The solver of Colella and Glaz (1985) is not exact in the strict sense because it involves a local linearization of the equation of state.

<!-- image -->

The use of simpler solvers in GR simulations is a concern because one-wave schemes behave significantly worse than full-wave solvers in the subsonic regime. The problem of excessive acoustic noise from the discontinuities introduced by the reconstruction is exacerbated because solvers like HLLE essentially have these discontinuities decay only into acoustic waves. This results in stronger numerical diffusion, but can also create artificial numerical noise because the diffusive terms in the HLLE or Rusanov flux can generate spurious pressure perturbations from isobaric conditions. While higher-order reconstruction can beat down the numerical diffusion for smooth flows, a strong degradation in accuracy is unavoidable in turbulent flow with structure on all scales.

Few attempts have as yet been made to quantify the impact of the more diffusive one-wave solvers in supernova simulations. In an idealized setup, the problem was addressed by Radice et al. (2015), who conducted simulations of stirred isotropic turbulence with solenoidal forcing with a turbulent Mach number of Ma /C24 0 : 35, two different Riemann solvers (HLLE vs. HLLC), and different reconstruction methods and grid resolutions. Even when using PPM reconstruction, they still found substantial differences between the HLLE and HLLC solver in the spectral properties of the turbulence with the HLLE runs requiring about 50 % more zones to achieve an equivalent resolution of the turbulent cascade to HLLC. Although it is not easy to extrapolate from the idealized setup of Radice et al. (2015) to corecollapse supernova simulations, one must clearly expect resolution requirements to depend sensitively on the Riemann solver. This is also illustrated by 2D supernova simulations comparing the HLLE and HLLC solver using the COCONUT-FMT code as shown in Fig. 4: starting from the initial seed perturbations, the HLLC model shows a faster growth of large-scale SASI shock oscillations during its linear phase and earlier emergence of parasitic instabilities (see Sect. 4.4 for the physics behind the SASI) due to the smaller amount of numerical dissipation. The evolution of the shock differs significantly during the first 100 ms of SASI activity, although the models become similar in terms of shock radius and shock asphericity later on. However, even then HLLC run consistently shows a higher entropy contrast and higher non-radial velocities within the gain region.

Reconstruction methods Similar concerns (not restricted to the low-Mach number regime) as for the simpler Riemann solvers can be raised about the order of the reconstruction scheme. There is certainly a clear divide between second-order piecewise linear reconstruction and higher-order methods like the PPM, WENO (weighted essentially non-oscillatory Shu 1997), and higher-order monotonicitypreserving (MP Suresh and Huynh 1997) schemes. In their simulations of forced subsonic turbulence, Radice et al. (2015) found similar differences between secondorder reconstruction using the monotonized central (MC; van Leer 1977) limiterone of the shaper second-order limiters-and runs using PPM or WENO as between the HLLC and HLLE solver. Again, the lower accuracy of second-order schemes is often clearly visible in full supernova simulations, which is again illustrated in Fig. 4 for the same setup as above. Similar to the HLLE run, the simulations using the MC limiter shows a delayed growth of the SASI and less small-scale structure.

<!-- image -->

<!-- image -->

x（km)

x（km)

Fig. 4 Snapshots of the entropy from 2D simulations of the 20 M /C12 progenitor of Woosley and Heger (2007) using the COCONUT-FMT code at post-bounce time of 137 ms (top row), 163 ms (middle row), and 226ms (bottom row). The left and rights halves of the panels in the left column show the results for the HLLE and HLLC solver with standard PPM reconstruction. The left and rights halves of the panels in the right column show the results for second-order reconstruction using the MC limiter and 6th-order extremum-preserving PPM reconstruction using the HLLC solver

Comparing the more modern higher-order reconstruction methods is much more difficult. For smooth problems like single-mode linear waves solutions, going beyond the original 4th-order PPM method of Colella and Woodward (1984) to methods of 5th order or higher can substantially reduce numerical dissipation; in the

<!-- image -->

optimal case, the dissipation decreases with the grid spacing D x as D x /C0 q for a q th order method (Rembiasz et al. 2017). However, the higher-order scaling of the numerical dissipation cannot be generalised to turbulent flow, because the dissipation of the shortest realizable modes at the grid scale does not increase as a higher power of q . Based on similarity arguments, one can work out that the effective Reynolds number of turbulent flow increases only as Re /C24 D x /C0 4 = 3 (Mu ¨ller 2016) and not as D x q as one might hope. The reason behind this limitation is that increasing the order of reconstruction does not increase the maximum wavenumber k max of modes that can be represented on the grid, it merely limits numerical dissipation to a narrow band of wavenumbers below k max.

For the moderately subsonic turbulent flow in core-collapse supernovae during the accretion phase, higher-order reconstruction often does not bring any tangible improvements for this reason. Figure 4 again shows this by comparing runs using standard 4th-order PPM and the 6th-order extremum-preserving PPM method of Colella and Sekora (2008). In both cases, the evolution of the shock is very similar, even though the phases of the SASI oscillations eventually falls out of sync. It is not obvious by visual inspection that the higher-order method allows smaller structures to develop. Only upon deeper analysis can small differences between the two methods be found, for example the model with extremum-preserving reconstruction maintains a measurably higher entropy contrast in the gain region and a slightly higher turbulent kinetic energy in the gain region.

There are nonetheless situations where it is useful to adopt extremum-preserving methods of very high order in global simulations of turbulent flow. First, such methods open up the regime of low Mach numbers to explicit Godunv-based codes. Using their APSARA code, Wongwathanarat et al. (2016) were able to solve the Gresho vortex problem (Gresho and Chan 1990) with little dissipation down to a Mach number of 10 /C0 4 with the extremum-preserving PPM method of Colella and Sekora (2008), which is about two orders of magnitude better than for the MC limiter (Miczek et al. 2015), and about one order of magnitude better than for standard PPM.

Modern higher-order methods can also be crucial in certain simulations of mixing at convective boundaries and nucleosynthesis. In the case of convective boundary mixing, this has been stressed and investigated by Woodward et al. (2010, 2014), who achieve higher accuracy for the advection of mass fractions in their PPMSTAR code by evolving moments of the concentration variables within each cell (which is somewhat reminiscent of the Discontinuous Galerkin method). They found that this Piecewise-Parabolic Boltzmann method only requires half the resolution of standard PPM to achieve the same accuracy (Woodward et al. 2010). Higher-order extremum-preserving methods may also prove particularly useful for minimizing the numerical diffusion of mass fractions in models of Rayleigh-Taylor mixing during the supernova explosion phase, but this is yet to be investigated.

## 2.1.3 High-Mach number flow

Some of the considerations for subsonic flow carry over to the supersonic and transsonic flow encountered during the supernova explosion phase where mixing

<!-- image -->

instabilities also lead to turbulence, but there are also problems specific to the supersonic regime.

Sonic points It is well known that the original Roe solver produces spurious expansion shocks in transsonic rarefaction fans, which needs to be remedied by some form of entropy fix (Laney 1998; Toro 2009). While other full-wave solvers-like the exact solver and HLLC-never fail as spectacularly as Roe's, they are still prone to mild instabilities at sonic points. Under adverse conditions, these instabilities can be amplified and turn into a serious numerical problem. In this case, it is advisable to switch to a more dissipative solver such as HLLE in the vicinity of the sonic point. In supernova simulations, this problem is sometimes encountered in the neutrino-driven wind that develops once accretion onto the PNS has ceased. It can also occur prior to shock revival in the infall region and severely affect the infall downstream of the instability, especially when nuclear burning is included. In this case, the problem can be easily overlooked or misidentified because it usually manifests itself as an unusually strong stationary burning front, which may seem perfectly physical at first glance.

Odd-even decoupling and the carbuncle phenomenon Full-wave solvers like the exact solver and HLLC are subject to an instability at shock fronts (Quirk 1994; Liou 2000): for grid-aligned shocks, insufficient dissipation in the direction parallel to the shock can cause odd-even decoupling in the solution, which manifests itself in artificial stripe-like patterns downstream of the shock (Fig. 5). When the shock is only locally tangential to a grid line, this instability can give rise to protrusions, which is known as the carbuncle phenomenon. In supernova simulations, odd-even decoupling was first recognized as a problem by Kifonidis et al. (2000), and since then the majority of supernova codes (e.g., PROMETHEUS, FLASH, COCONUT, FORNAX) have opted to handle this problem by adaptively switching to the more dissipative HLLE solver at strong shocks following the suggestion of Quirk (1994). The CHIMERA code (Bruenn et al. 2018) adopts the alternative approach of a local oscillation filter (Sutherland et al. 2003), which has the advantage of not degrading the resolution in the direction perpendicular to the shock, but has the drawback of allowing the instability to grow to a minute level (which may be undetectable in

<!-- image -->

Fig. 5 Odd-even decoupling in a 2D core-collapse supernova simulation of a 20 M /C12 star with the COCONUT that uses the HLLC solver everywhere instead of switching to HLLE at shocks. The left and right panels show the radial velocity in units of the speed of light and the entropy s in units of k b = nucleon about 10 ms after bounce. The characteristic radial streaks from odd-even decoupling are clearly visible behind the shock

<!-- image -->

practice) before smoothing is applied to the solution. The carbuncle phenomenon can also occur in Richtmyer-type artificial viscosity schemes and be cured by modifying the artificial viscosity (Iwakami et al. 2008). The carbuncle instability remains a subject of active research in computational fluid dynamics, and a number of papers (e.g., Nishikawa and Kitamura 2008; Huang et al. 2011; Rodionov 2017; Simon and Mandal 2019) have attempted to construct Riemann solvers or artificial viscosity schemes that avoid the instability without sacrificing accuracy away from shocks, and may eventually prove useful for supernova simulations.

Kinetically-dominated flow In HRSC codes that solve the total energy equation, one obtains the mass-specific internal energy e by subtracting the kinetic energy v 2 = 2 from the total energy e . In high Mach-number flow, one has e /C29 e and v 2 = 2 /C29 e , and hence subtracting these two large terms can introduce large errors in the internal energy density and the pressure and sometimes leads to severe stability problems. A similar problem can occur in magnetically-dominated regions in MHD. Sometimes the resulting stability problems can be remedied by evolving the internal energy equation

<!-- formula-not-decoded -->

instead of Eq. (3) in regions of high Mach number or low plasmab . However, in doing so one sacrifices strict energy conservation, and hence one should apply this recipe as parsimoniously as possible.

## 2.2 Treatment of gravity

Convective burning and core-collapse supernovae introduce specific challenges in the treatment of gravity. In the subsonic flow regimes, one needs to be wary of introducing undue artifical perturbations from hydrostatic equilibrium and take care to avoid secular conservation errors. Moreover, in the core-collapse supernova problem, general relativistic effects become important in the vicinity of the PNS.

## 2.2.1 Hydrostatic balance and conservation properties

For nearly hydrostatic flow, one has r P /C25 /C0 q r U , but this near cancellation is not automatically reflected in the numerical solution when using a Godunov-based scheme. Instead, the stationary numerical solution may be one with non-zero advection terms that are exactly (but incorrectly) balanced by the gravitational source term (Greenberg and Leroux 1996; LeVeque 1998a). Schemes that avoid this pathology are called well-balanced . The proper cancellation between the pressure gradient and the gravitational source term is particularly delicate if those two terms are treated in operator-split steps. Different methods have been proposed to incoroprate well-balancing into Godunov-based schemes. One approach is to use piecewise hydrostatic reconstruction (e.g., Kastaun 2006; Ka ¨ppeli and Mishra 2016). A related technique suggested by LeVeque (1998a) introduces discrete jumps

<!-- image -->

in the middle of cells to obtain modified interface states for the Riemann problem and absorb the source terms altogether.

In practice, these special techniques are not used widely in the field for two reasons. First, it is not trivial to general these schemes to achieve higher-order accuracy. Second, one already obtains a very well-balanced scheme by combing higher-order reconstruction, an accurate Riemann solver, and unsplit time integration. For split schemes, one can ensure a quite accurate cancellation of the pressure gradient and the source term by including a characteristic state correction as described by Colella and Woodward (1984) for the original PPM method.

Nevertheless, the cancellation of the pressure gradient and the gravitational source term in hydrostatic equilibrium is usually not perfect and typically leads to minute odd-even noise in the velocity field that is almost undetectable by eye. Computing the gravitational source term q v /C1r U in the energy equation using such a noisy velocity field v can lead to an appreciable secular drift of the total energy. For example, spurious energy generation can stop proto-neutron star cooling on simulation time scales longer than a second (Mu ¨ller 2009). This problem can be circumvented by discretizing the energy equation starting from the form (Mu ¨ller et al. 2010)

<!-- formula-not-decoded -->

This guarantees exact total energy conservation if the time derivative of the gravitational potential is zero. Under certain conditions, exact total energy conservation can be achieved for a time-dependent self-gravitating configuration as well, and the method can also be generalized to the relativistic case.

In principle, one can also implement the gravitational source term (in the Newtonian approximation) in the momentum equation in a conservative form by writing q g as the divergence of a gravitational stress tensor (Shu 1992). Such a scheme has been implemented by Livne et al. (2004) in the VULCAN code. However, this procedure involves a more delicate modification of the equations than in case of the energy source term, because it essentially amounts to replacing q by the finitedifference representation of the Laplacian ð 4 p G Þ /C0 1 DU in the momentum source term. Unless the solution for the gravitational potential is extremely accurate, large acceleration errors may thus arise. Moreover, this approach does not work for effective relativistic potentials (see Sect. 2.2.2). For these reasons, the conservative form of the gravitational source term has not been used in practice in other codes. Even though the issue of momentum conservation is of relevance in the context of neutron star kicks, conservation errors do not seem to affect supernova simulation results qualitatively in practice, and post-processing techniques can be used to infer neutron star velocities from simulations with good accuracy (Scheck et al. 2006).

## 2.2.2 Treatment of general relativity

In core-collapse supernova simulations, the relativistic compaction of the protoneutron star reaches GM = Rc 2 ¼ 0 : 1-0.2 even for a normal PNS mass M and a

<!-- image -->

somewhat extended radius R of the warm PNS. Infall velocities of 0.15-0.3 c are encountered. Hence general relativistic (GR) and special relativistic effects are no longer negligible, though the latter is more critical for the treatment of the neutrino transport than for the hydrodynamics. For very massive neutron stars, cases of black hole formation, or jet-driven explosions, relativistic effects can be more pronounced.

A variety of approaches is used in supernova modelling to deal with relativistic effects. Purely Newtonian models have now largely been superseded. Using Newtonian gravity results in unphysically large PNS radii, and, as a consequence, lower neutrino luminosities and mean energies and worse heating heating conditions than in the relativistic case, even though the stalled accretion shock radius is larger than in GR before explosion (Mu ¨ller et al. 2012b; Kuroda et al. 2012; Lentz et al. 2012; O'Connor and Couch 2018b). As an economical alternative, one can retain the framework of Newtonian hydrodynamics but incoroprate relativistic corrections in the gravitational potential based on the TOV equation (Rampp and Janka 2002). This approach was subsequently refined by Marek et al. (2005), Mu ¨ller et al. (2008) to account for some inconsistencies between the use of Newtonian hydrodynamics and a potential based on a relativstic stellar structure equation, but full consistency can never be achieved in the pseudo-Newtonian approach. In the multi-D case, the relativistic potential replaces the monopole of the Newtonian potential, while higher multipoles are left unchanged. From a purist point of view, this pseudo-Newtonian approach is delicate because one sacrifices global conservation laws for energy and momentum (which would still hold in a more complicated form in an asymptotically flat space in full GR). In practice, this is less critical; in PNS cooling simulations by Hu ¨depohl et al. (2010) the total emitted neutrino energy was found to agree with the neutron star binding energy (computed from the correct TOV solution) to within 1 % for the modified TOV potential (Case A) of Marek et al. (2005).

If the framework of Newtonian hydrodyanmics is abandoned, one may still opt for an approximate method to solve for the space-time metric as in the COCONUT code (Dimmelmeier et al. 2005; Mu ¨ller et al. 2010; Mu ¨ller and Janka 2015). Elliptic formulations such as CFC (conformal flatness conditions, Isenberg 1978) and xCFC (a modification of CFC for improved numerical stability; Cordero-Carrio ´n et al. 2009) can be cheaper and more stable than free-evolution schemes based on the 3 ? 1 decomposition (for reviews of these techniques, see Baumgarte and Shapiro 2010; Lehner and Pretorius 2014) and maximally constrained schemes (Bonazzola et al. 2004; Cordero-Carrio ´n et al. 2012). However, full GR supernova simulations without the CFC approximation and with multi-group transport have also become possible recently (Roberts et al. 2016; Ott et al. 2018; Kuroda et al. 2016b, 2018). Although CFC remains an approximation, it is exact in spherical symmetry, and comparisons with free-evolution schemes have shown excellent agreement in the context of rotational collapse have shown excellent agreement even for rapidly spinning progenitors (Ott et al. 2007).

Comparisons of pseudo-Newtonian and GR simulations have demonstrated that using an effective potential is at least sufficient to reproduce the PNS contraction, the shock evolution, and the neutrino emission in GR very well (Liebendo ¨rfer et al. 2005; Mu ¨ller et al. 2010, 2012b). While Mu ¨ller et al. (2012b) still found better

<!-- image -->

heating conditions in the GR case than with an effective potential in their 2D models, this comparison was not fully controlled in the sense that two different hydro solvers were used, and the effect was related to subtle differences in the PNS convection zone, which may well be related to factors other than the GR treatment (cf. Sect. 4.7). Further code comparisons are desirable to resolve this. The pseudoNewtonian approach, does, however, systematically distort the eigenfrequencies of neutron star oscillation modes (Mu ¨ller et al. 2008). In particular, the frequency of the dominant f-/g-mode seen in the gravitational wave spectrum is shifted up by 15-20% compared to the correct relativistic value (Mu ¨ller et al. 2013).

## 2.2.3 Poisson solvers

In the Newtonian approximation, the gravitational field is obtained by solving the Poisson equation

<!-- formula-not-decoded -->

In constrained formulations of the Einstein equations like (x)CFC, one encounters non-linear Poisson equations.

In simulations of supernovae and the late convective burning stages, the density field usually only deviates modestly from spherical symmetry and is not exceedingly clumpy (except in the case of mixing instabilities in the envelope during the explosion phase when self-gravity is less important to begin with). For this reason, the usual method of choice for solving the Poisson equation (even in Cartesian geometry) is to use the multipole expansion of the Green's function (Mu ¨ller and Steinmetz 1995). Typically, no more than 10-20 multipoles are needed for good accuracy, and very often only the monopole component is retained. Other methods have been used occasionally, though, such as pseudospectral methods (Dimmelmeier et al. 2005) and finite-difference solvers (e.g., Burrows et al. 2007b), and the FFT (Hockney 1965; Eastwood and Brownrigg 1979) is a viable option for Cartesian simulations.

Although it yields accurate results at fairly cheap cost, some subtle issues can arise with the multipole expansion. When projecting the source density onto spherical harmonics Y ' m to obtain multipole components ^ q ' m ¼ R Y /C3 ' m q d X , a naive step-function integration can lead to a self-potential error (Couch et al. 2013) and destroy convergence with increasing mulitpole number N ' . This can be avoided either by performing the integrals over spherical harmonics Y /C3 ' m analytically (Mu ¨ller 1994), or by using a staggered grid for the potential (Couch et al. 2013). The accuracy of the solution can also be degraded if the central mass concentration moves away from the center of the grid, which can be cured by off-centering the multipole expansion (Couch et al. 2013). Problems with off-centred or clumpy mass distributions can be cured completely if an exact solver is used. In Cartesian geometry, this can be accomplished econmically using the FFT, and an exact solver for spherical polar grid using a discrete eigenfunction expansion has recently been developed as well (Mu ¨ller and Chan 2019). On spherical multi-patch grids, the efficient parallelization and computation of integration weights requires some

<!-- image -->

thought and has been addressed by Almansto ¨tter et al. (2018) and Wongwathanarat (2019).

## 2.3 Reactive flow

Nuclear burning is the principal driver of the flow for core and shell convection in the late, neutrino-cooled evolutionary stages of supernova progenitors. In corecollapse supernovae nuclear dissociation and recombination play a critical role for the dynamics and energetics, and one of the key observables, the mass of 56 Ni, is determined by nuclear burning.

Approaches to nuclear transmutations differ widely between simulation codes, and range from the assumption of nuclear statistical equilbrium (NSE) everywhere in some core-collapse supernova models to rather large reaction networks. Naturally,the appropriate level of sophistication depends on the regime and the observables of interest. The theory of nuclear reaction networks is too vast to cover in detail here, and we can only touch a few salient points related to their integration into hydrodynamics codes. For a more extensive coverage, we refer to textbooks and reviews on the subject (Clayton 1968; Arnett 1996; Mu ¨ller 1998; Timmes 1999; Hix and Meyer 2006; Iliadis 2007).

Burning regimes As stellar evolution proceeds towards collapse, the ratio of the nuclear time scale to both the sound crossing time scale and convective time scale decreases, and the nuclear reaction flow involves an increasing number of reactions. The burning of C, Ne, and to some extent of O is dominated by an overseeable number of main reaction channels, and the relevant reaction rates are slow compared to the relevant hydrodynamical time scales. During oxygen burning, quasi-equilibrium clusters begin to appear and eventually merge into one or two big clusters during Si burning (Bodansky et al. 1968; Woosley et al. 1973) that are linked by slow ''bottleneck'' reactions. For sufficiently high temperatures, NSE is established and the composition only depends on density q , temperature T , and the electron fraction Y e and is given the Saha equation. At higher densities during corecollapse, the assumption of non-interacting nuclei break down, and a high-density equation of state is required (see Lattimer 2012; Oertel et al. 2017; Fischer et al. 2017, for recent reviews); this regime is not of concern here because the flow can be treated as non-reactive.

Simple approaches In core-collapse supernova simulations, one sometimes simply assumes NSE everywhere, which amounts to an implicit release of energy at the start of a simulation. Although the Si and O shell will still collapse in the wake of the Fe core, this is somewhat problematic, especially for long-time simulations where the effect on the infall is bound to be more pronounced. For mitigating potential artifacts from the inconsistency of the composition and equation of state with the underlying stellar evolution model, it can be useful to initialise supernova simulations using the pressure rather than the temperature of progenitor model.

A considerably better and very cheap approach, known as ''flashing'', is to use a few key a -elements and non-symmetric iron group nuclei in addition to protons,

<!-- image -->

neutrons, and a -particles and burn them instantly into their reaction products and eventually into NSE upon reaching certain threshold temperatures (Rampp and Janka 2002). Such an approach can capture the energetics of explosive burning in the shock and the freeze-out from NSE in neutrino-driven outflows reasonably well, but only gives indicative results on the composition of the ejected matter. The choice of the proper NSE threshold temperature T NSE J 5GK can be particularly delicate, since this depends on the entropy and expansion time scale of the outflow and can critically affect the degree of recombination in the neutrino-heated ejecta.

For simulations of convective burning, a smooth behavior of the nuclear source terms is required, but for C, Ne, or O burning, one can still resort to simple fit formulae and only track the composition of the main fuel and ash if the goal is merely to understand the dynamics of the flow (e.g., Kuhlen et al. 2003; Jones et al. 2017; Cristini et al. 2016, 2017; Andrassy et al. 2020). Often these source terms (or also rates in calculations with veritable network) are rescaled if low convective Mach numbers make simulations in the physical regime unfeasible. This can be useful for exploring the parameter dependence of flow phenomena, but caution is required because safe extrapolation to the physical regime may also require a rescaling of other terms (e.g., thermal diffusivity, neutrino losses).

Reaction networks In 2D, Baza ´n and Arnett (1997) already conducted simulations of convective burning with 123 species, but the use of large networks in 3D simulations is still prohibitively expensive. Modern 3D simulations of convective burning with the PROMPI (e.g., Meakin and Arnett 2006; Moca ´k et al. 2018), PROMETHEUS (Mu ¨ller et al. 2016b; Yadav et al. 2020), FLASH (Couch et al. 2015) and 3DNSNE (Yoshida et al. 2019) codes have therefore only use networks of 19-25 species consisting of a -elements, light particles, and at most a few extra iron-group elements. In multi-D supernova simulations with neutrino transport the use of such networks is feasible (von Groote 2014; Bruenn et al. 2013, 2016; Wongwathanarat et al. 2017), though they have not been used widely yet. It is critical that such reduced reaction networks appropriately account for side chains and the effective reaction flow between light particles (Weaver et al. 1978; Timmes et al. 2000). Their use is problematic for Si burning which requires networks of more than a hundred species to accurately capture the quasi-equilibrium clusters and the effects of deleptonization (Weaver et al. 1978) and for freeze-out from NSE with considerable neutron excess. Larger networks or special methods for quasiequilibrium (Weaver et al. 1978; Hix et al. 2007; Guidry et al. 2013) will be required for reliable multi-D simulations of convective Si burning.

Coupling to the hydrodynamics Some numerical issues arise when a nuclear network is coupled to a Eulerian hydrodynamics solver, or even if the composition is just tracked as a passive tracer. One such problem concerns the conservation of partial masses, which is guaranteed analytically by a conservation Eq. (4) for each species i ,

<!-- image -->

<!-- formula-not-decoded -->

This equation can be solved using standard, higher-order finite-volume techniques. However, the solution also has to obey the constraint

X

<!-- formula-not-decoded -->

which is not fulfilled automatically by the numerical solution, unless flat reconstruction for the mass fractions is employed. One could enforce this constraint by rescaling the mass fractions to sum up to unity, but this would violate the conservation of partial masses. Plewa and Mu ¨ller (1999) developed the Consistent Multifluid Advection (CMA) method as the standard treatment to ensure both minimal numerical diffusion of mass fractions and enforce conservation of partial masses. This method involves a rescaling and coupling of the interpolated interface values of the various mass fractions. Plewa and Mu ¨ller (1999) demonstrated that simple methods for the advection of mass fractions can easily result in wrong yields by a factor of a few for some isotopes in supernova explosions.

Another class of problems is related to advection errors and numerical diffusion, especially at contact discontinuities and shocks, which can lead to artificial detonations or an incorrect propagation of physical detonations (Colella et al. 1986; Fryxell et al. 1989; Mu ¨ller 1998). To ensure that detonations propagate at the correct physical velocity, nuclear burning should be switched off in shocked zones (Mu ¨ller 1998). 5 Due to the extreme temperature dependence of nuclear reaction rates, similar problems can arise away from discontinuities due to advection errors that produce a small level of noise in the temperature. Artificial detonations can easily develop in highly degenerate regions and around sonic points. Eliminating such artifacts may require appropriate switches for pathological zones or very high spatial resolution (e.g., Kitaura et al. 2006).

## 3 Late-stage convective burning in supernova progenitors

In the Introduction, we already outlined the motivation for multi-D simulations of supernova progenitors in broad terms. On the most basic level, multi-D models are needed to properly intialize supernova simulations and provide physically correct seed perturbations for the instabilities that develop after collapse and in the explosion phase. This does not, in fact, presuppose that 1D stellar evolution models incorrectly predict the overall spherical structure of pre-supernova progenitors; in principle such an initialization might involve nothing but adding some degrees of freedom to 1D stellar evolution models without any noticeable change of the spherically averaged stratification. Historically, however, simulations of late-stage convection have focused on deviations of the multi-D flow from the predictions of traditional mixing-length theory (MLT; Biermann 1932; Bo ¨hm-Vitense 1958; Weiss

5 Note also the use of front-tracking methods for unresolvable burning fronts (Reinecke et al. 2002; Leung and Nomoto 2019), which are commonly used for modelling Type Ia supernovae and the O deflagration in electron-capture supernova progenitors.

<!-- image -->

et al. 2004) and not evolved progenitor models up to core collapse, whereas the initialization problem has only been tackled recently by Couch et al. (2015), Mu ¨ller et al. (2016b, 2019), Yadav et al. (2020). In this section, we therefore address the interior flow in convective regions and boundary effects first before specifically discussing multi-D pre-supernova models.

## 3.1 Interior flow

Let us first consider the flow within convectively unstable regions. In MLT as implemented in modern stellar evolution codes such as KEPLER (Weaver et al. 1978; Heger and Woosley 2010) and MESA (Paxton et al. 2011), the convective velocity v conv in such regions is tied to the superadiabaticity of the density gradient as encoded by the Brunt-Va ¨isa ¨la ¨ x BV frequency and the local pressure scale height K ,

p

ffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffi ffi

<!-- formula-not-decoded -->

where a is a tuneable parameter of order unity, and the MLT density contrast dq is obtained from the the difference between the actual and density gradient o q = o r and the adiabatic density gradient ð o q = o P Þ s ð o P = o r Þ ,

/C20

/C18

/C19

/C21

/C18

/C19

<!-- formula-not-decoded -->

Note that stellar evolution textbooks usually express the convective velocity and density contrast in terms of the difference between the actual and adiabatic temperature gradient (Clayton 1968; Weiss et al. 2004; Kippenhahn et al. 2012), but Eqs. (30) and (31) are fully equivalent formulations that often prove less cumbersome.

Using Eq. (30) for the convective velocity, Eq. (31) for the MLT density contrast, and the temperature contrast d T ¼ ð o T = o ln q Þ P ð dq = q Þ , we then obtain the convective energy flux F conv (Kippenhahn et al. 2012; Weiss et al. 2004),

<!-- formula-not-decoded -->

where cP is the specific heat at constant pressure, and a e is another tunable nondimensional parameter. Similarly, by estimating the composition contrast d Xi using the local gradient as d Xi ¼ a X K o Xi = o r , we obtain the partial mass flux for species i

<!-- formula-not-decoded -->

where a X is again a dimensionless parameter. When comparing 1D stellar evolution models to each other or to multi-D simulations, one must bear in mind that slightly different normalization conventions for Eqs. (32) and (33) are in use. Regardless of these ambiguities, these coefficients are of order unity, for example the KEPLER code

<!-- image -->

uses aa e ¼ 1 = 2 and aa X ¼ 1 = 6 (Weaver et al. 1978), which can be conveniently interpreted as a ¼ 1, a e ¼ 1 = 2 and a X ¼ 1 = 6 (Mu ¨ller et al. 2016b),

In order to connect more easily to multi-D simulations, it is useful restate the assumptions and consequences of MLT (without radiative diffusion) in a slightly different language. Equation (30) can also be written as

<!-- formula-not-decoded -->

which we can interpret as a balance between the rate of buoyant energy generation ( a 2 dq = q v conv) and turbulent dissipation ( /C15 /C24 v 3 conv = K ). Furthermore the work done by bouyancy must ultimately be supplied by nuclear burning. Using thermodynamic relations, we find that the potential energy Kdq = q g liberated by bubbles rising or sinking by one mixing length is of the order of the enthalpy contrast d h of the bubbles, which roughly equals the integral of the nuclear energy generation rate \_ q nuc over one turnover time s ¼ K = v conv ,

<!-- formula-not-decoded -->

Together, Eqs. (34) and (35) lead to a scaling law v conv /C24 ð \_ q nuc K Þ 1 = 3 for the typical value of v conv in a convective shell.

In nature, balance between nuclear energy generation, buoyant energy generation, and turbulent dissipation is usually established over a few turnover times. On longer time scales, active burning shells also adjust by expansion or contraction until the total nuclear energy generation rate and neutrino cooling rate balance each other (Woosley et al. 1972), with the nuclear burning dominating in the inner region and neutrino cooling dominating in the outer region of the shell (Fig. 6).

<!-- image -->

Fig. 6 Energy generation rate \_ q nuc (black), neutrino cooling rate j \_ q m j (red), and convective velocity v conv (blue) in an 18 : 88 M /C12 progenitor (1D model, discussed in Yadav et al. (2020) in the innermost shells about 1hr before collapse. At this stage, balanced power still obtains. Nuclear energy generation dominates at the bottom of the shells, while neutrino cooling dominates in the outer layers. The integrated energy generation and cooling rate for the entire shell, which are given by the areas under the black and red curve, nearly balance each other

<!-- image -->

Because of the extremely strong temperature sensitivity of the burning rates, this state of balanced power is difficult to maintain when setting up multi-D simulations and will only be reestablished over a long, thermal time scale. 6 In fact, the problem of thermal adjustment has not yet been rigorously analyzed for any multi-D model yet, and insufficient simulation time for thermal adjustment is a concern that needs to be addressed in future. However, the problem of thermal adjustment is mitigated during the latest phases of shell convection prior to collapse: As the core and the surrounding shells contract, the nuclear burning rates accelerate to a point where neutrino cooling and shell expansion by P d V work can no longer re-establish thermal balance on the contraction time scale of the core, and the state of balanced power is physically broken.

Two-dimensional simulations of convective burning The first attempts to simulate late-stage convection in massive stars by Arnett (1994), Bazan and Arnett (1994, 1998) targeted oxygen burning in a 20 M /C12 star in a 2D shellular domain with the PROMETHEUS code using a small, 12-species reaction network and neutrino cooling by the pair process. Starting from a simulation on a small wedge of 18 /C14 in Arnett (1994), Bazan and Arnett (1994, 1998) subsequently considered broader wedges in cylindrical symmetry of up to 135 /C14 with a resolution of up to 460 /C2 128 grid cells, as well as cases with meridional symmetry on a 2D grid ð r ; u Þ in radius and longitude. The simulations were invariably limited to short periods (up to 400 s in Bazan and Arnett 1998) and only a few convective turnover times. One simulation (Baza ´n and Arnett 1997) also tackled Si burning in 2D with a large network of 123 nuclei. These first-generation 2D models invariably found violent convective motions with Mach numbers of 0.1-0.2 and velocities about an order of magnitude above the MLT predictions, which cannot be accounted for by the aforementioned ambiguities in the definition of the dimensionless coefficients. The convective structures invariably tended to grow to the largest angular scale allowed by the chosen wedge geometry, and large density perturbations were found at the convective boundaries. Bazan and Arnett (1998) also stressed the high temporal variability of the convective flow, going so far as to question whether a steady state is ever established before collapse. Longer simulations of the same 20 M /C12 model over 1200 ms by the same group using the VULCAN code of (Livne 1993) showed the emergence of a steady state, albeit quite different from the 1D stellar evolution model due to convective boundary mixing (Asida and Arnett 2000).

To a large extent, the pronounced differences between these first-generation simulations and MLT predictions stem from the assumption of 2D flow. In 2D turbulence, the energy cascade is artificially inverted and goes from small to large scales (Kraichnan 1967). As a result, the flow tends to organise itself into large vortices, and dissipation occurs primarily in boundary layers (Falkovich et al. 2017; Clercx and van Heijst 2017).

6 Note that this thermal time scale is more difficult to define than during early burning stages where radiative diffusion is important.

<!-- image -->

Three-dimensional simulations of convective burning Consequently, 3D simulations of convective burning obtained considerably smaller convective velocities. The first 3D, full 4 p solid angle models of O shell convection (along with models of core hydrogen burning) were presented by Kuhlen et al. (2003) for a 25 M /C12 star with and without rotation. Their simulations used the anelastic pseudospectral code of Glatzmaier (1984) to follow convection for about 90 turnovers in the non-rotating case, and approximated the burning and neutrino cooling rates by power-law fits. Different from the earlier 2D models, they found convective velocities in good agreement with the 1D MLT prediction in the underlying stellar evolution model, but the still observed the emergence of large-scale flow patterns.

The use of simplified burning and neutrino loss rates, the anelastic approximation, and an explicit turbulent diffusivity in Kuhlen et al. (2003) still posed a concern, which was subsequently addressed by a series of 2D and 3D simulations of O and C burning (Meakin and Arnett 2006, 2007a, b; Arnett et al. 2009) in wedgeshaped domains using the compressible PROMPI code and a larger reaction network (25 species) than in the first generation of 2D models. These simulations confirmed the significantly less violent nature of 3D convection compared to 2D (Meakin and Arnett 2006, 2007b), and established good agreement between elastic and anelastic simulations on the convective velocities and fluctuations of thermodynamic quantities in the interior of convective zones, though anelastic codes cannot model fluctuations at convective boundaries very well (Meakin and Arnett 2007a). They also found balance between buoyant driving and turbulent dissipation (which is essentially a restatement of the basic assumption of MLT) and observed rough equipartiton between the radial and non-radial contributions to the turbulent kinetic energy (Arnett et al. 2009). Their models still revealed differences from MLT in detail, such as different correlation lengths for velocity and temperature and a nonvanishing kinetic energy flux (Meakin and Arnett 2007b). Moreover, Meakin and Arnett (2010) suggested that the implicit identification of the pressure scale height with the dissipation length in MLT might lead to an underestimation of the convective velocities. More recent work by the same group has stressed the time variability of the convective flow (Arnett and Meakin 2011a, b) and criticized the MLT assumption of quasi-stationary convective velocities. Specifically Arnett and Meakin (2011b) pointed to strong fluctuations in the turbulent kinetic energy in the 3D oxygen shell burning simulation in a 23 M /C12 star by Meakin and Arnett (2007b), which they attempted to motivate by recourse to the Lorenz model for convection in the Boussinesq approximation. The connection between the simulations of convective burning and the Lorenz model for a viscous-conductive convection problem remains rather opaque, however.

More recent work on 3D convection by other groups has vindicated rather than undermined MLT as an approximation for the interior of convective zones. Mu ¨ller et al. (2016b) conducted a 4 p -simulation of O burning in an 18 M /C12 star up to the point of collapse and found that convection reaches a quasi-stationary state after a few turnovers with only small fluctuations in the turbulent kinetic energy. In line with MLT and as in Arnett et al. (2009), the average convective velocity is well described by a balance of turbulent dissipation and buoyant driving in their model, and is in turn related to the average nuclear energy generation rate \_ q nuc as

<!-- image -->

<!-- formula-not-decoded -->

and even the profiles of the radial component of the turbulent velocity perturbation are in good agreement with the corresponding 1D stellar evolution model. A similar scaling was reported by Jones et al. (2017) based on idealized high-resolution simulations of O burning with a simple EoS and parameterized nuclear source terms and by Cristini et al. (2017) based on simulations of C burning in planar 3D geometry, also with a parameterized (and artificially boosted) nuclear source term. Jones et al. (2017) verified this scaling over a wider range of convective luminosities and Mach numbers by applying different boost factors to the nuclear generation rate.

Regarding the dominant scales of the convective flow, the recent global 3D shell burning simulations (Chatzopoulos et al. 2014; Couch et al. 2015; Mu ¨ller et al. 2016b; Jones et al. 2016; Yadav et al. 2020) confirm the emergence of large eddies with low angular wavenumber ' that stretch across the entire convective zone. Mu ¨ller et al. (2016b) verified quantitatively that the peak of the turbulent energy spectrum in ' agrees well with the wavenumber of the first unstable convective mode at the critical Rayleigh number (Chandrasekhar 1961; Foglizzo et al. 2006),

<!-- formula-not-decoded -->

Further simulations that also explored thinner shells (Mu ¨ller et al. 2019) show a shift towards higher ' and corroborate this scaling as illustrated in Fig. 7. Beyond this dominant wavenumber, the turbulence exhibits a Kolmogorov spectrum (Chatzopoulos et al. 2014; Mu ¨ller et al. 2016b).

Naturally, the modern 3D models still exhibit differences to MLT in detail even within convective zones. For example, x 2 BV often changes sign in the outer parts of a convective layer in 3D, indicating that the spherically-averaged stratification is nominally stable (Moca ´k et al. 2009; Mu ¨ller et al. 2016b). Mu ¨ller et al. (2016b) also remark that the spherically-averaged mass fraction profiles tend to be flatter in 3D than in 1D, due to the usual asymmetric choice a X ¼ a e = 3 for the MLT parameters for material diffusion and energy transport, which probably ought to be replaced by a X ¼ a e . A rigorous approach to quantify the structure of the convective flow and the differences between 3D and 1D models is available in the form of spherical Reynolds decomposition, which has been pursued systematically by Viallet et al. (2013), Moca ´k et al. (2014) and Arnett et al. (2015). The mere form of the Reynolds-averaged equations for bulk (i.e., spherically-averaged) and fluctuating quantities dictates that such an analysis invariably finds dozens of terms that are implicitly set to zero in MLT.

Assessment How are we to evaluate these commonalities and differences between 3D simulations and 1D stellar evolution flow? For most purposes, the question is not whether effects are missing in MLT-based 1D models (since the very purpose of an approximation like MLT is to retain only the leading effects), but whether those missing effects matter over secular time scales or have an impact during the supernova explosion. As we shall discuss in detail in Sect. 4.5, the presence of

<!-- image -->

3

Fig. 7 Dependence of the dominant eddy scale on the shell geometry illustrated by slices through 3D supernova progenitor models with convective burning and their turbulent energy spectra. The 2D slices show the radial velocity at the onset of collapse in progenitors of 12 M /C12 with metallicity Z ¼ 0 (top right), 14 : 07 M /C12 with Z ¼ Z /C12 (top left) , and 12 : 5 M /C12 with Z ¼ Z /C12 (bottom left) with active convective O shells. The bottom right panel shows turbulent energy spectra E ð ' Þ computed from the radial velocity around the center of the convective zone. The dominant wavenumber expected from Eq. (37) is indicated at the top; note that there is an uncertainty because the outer boundaries of the convective zones are fuzzy. The dotted lines show the slope of a Kolmogorov spectrum. The plots for the 12 M /C12 and 12 : 5 M /C12 models have been adapted from Mu ¨ller et al. (2019). Image reproduced with permission, copyright by the authors

<!-- image -->

asymmetries in convective shells indeed matters during the supernova, but the fact is also that MLT and linear perturbation theory appear to predict the relevant parameters-the velocities and dominant scales of convective eddies-quite well. As far as the secular evolution of convective burning shells is concerned, there is little evidence that MLT does not adequately describe the flow within convective shells. There is typically good agreement in critical parameters for the shell evolution like the total nuclear burning rate. Many effects that MLT captures inaccurately and matter critically in models of convective envelopes and stellar atmospheres-such as the precise deviation of the stratification from superadiabticity-are of minor importance for the bulk evolution of massive stars during the late burning stages. For more tangible consequence of multi-D effects on secular time scales, we need to consider convective boundaries in Sect. 3.3.

<!-- image -->

## 3.2 Supernova progenitor models

Simulations to the presupernova stage Only a few models of convective burning have yet been carried up to the point of core collapse (Couch et al. 2015; Mu ¨ller 2016; Mu ¨ller et al. 2016b, 2019; Yadav et al. 2020; Yoshida et al. 2019) because of several obstacles. In order to accurately follow the composition changes and the deleptonization in the Fe core and Si shell (i.e., in the NSE and QSE regime) that drive the evolution towards collapse, reaction networks with well over a hundred nuclei are required (Weaver et al. 1978). This is feasible in principle, but yet impractical for well-resolved 3D simulations up to collapse. Furthermore, the initial transient phase and imperfect hydrostatic equilibrium after the mapping from 1D to multi-D may artificially delay the collapse.

Two different strategies have been employed to circumvent these problems. In their simulation of Si burning in a 15 M /C12 star for 160 s, Couch et al. (2015) used an extended 21-species a -network with some iron group nuclei added to model core deleptonization (Paxton et al. 2011). In order to force the core to collapse, they increased the electron capture rate on 56 Fe by a factor of 50, and their 3D model in fact reaches collapse more than six times faster than the corresponding 1D stellar evolution model. This approach is problematic because any modification of the contraction time scale of the core also affects the burning in the outer shells (Mu ¨ller et al. 2016b). Using the same 21-species network, Yoshida et al. (2019) managed to evolve a 3D simulation of a 25 M /C12 star and several 2D simulations of different progenitors for the last /C24 100 s without modifying the deleptonization rate. This suggests that multi-D models can be evolved somewhat self-consistently to collapse even though the short simulation time is a concern in this particular case, since it remains unclear to what extent the results are affected by the initial transient.

The 3D studies of O shell burning in various progenitors by Mu ¨ller (2016), Mu ¨ller et al. (2016b), Mu ¨ller et al. (2019) and Yadav et al. (2020) have followed a different approach and circumvented the problems of QSE and deleptonization by excising the major part of the Fe core and Si shell and replacing them with an inner boundary condition that is contracted according to a mass shell trajectory from the corresponding 1D stellar evolution model. This approach can be justified for many progenitors, which have no active convective Si shell, or only weak convection in the Si shell.

Evolution towards collapse The convective flow in the contracting burning shells shortly before collapse exhibits few noteworthy differences to the burning in quasihydrostatic shells described in Sect. 3.1. The 3D simulations of the different groups (Couch et al. 2015; Mu ¨ller et al. 2016b; Yoshida et al. 2019) all show the emergence of modes with a dominant wavelength of the order of the shell width according to Eq. (37), and as far as comparisons have been performed, the convective velocities remain in good agreement with MLT until shortly before collapse. It is noteworthy, however, that the convective velocities and Mach numbers tend to increase significantly during the last minutes before collapse because the temperature at the base of the inner shells, and hence the burning rate, increase as they contract in the wake of the core. The convective velocities then

<!-- image -->

freeze out shortly before collapse once the burning rate changes on a time scale shorter than the turnover time scale. This freeze-out seems to be captured adequately by time-dependent MLT so that 1D stellar evolution models provide good estimates for the convective velocities at the onset of collapse. Bigger differences between 1D and 3D progenitor models can occur in case of small buoyancy barriers between the O, Ne, and C shell, in which case 3D models are more likely to undergo a shell merger (Yadav et al. 2020).

The evolution of the convective shells during collapse will be discussed in Sect. 4.5.

Progenitor dependence Since 3D simulations indicate that convective velocities and eddy scales can be estimated fairly well from 1D stellar evolution models, one can already roughly outline the progenitor dependence of convective shell properties as shown by Collins et al. (2018). Considering the active Si and O shell burning shells at the onset of collapse in over 2000 progenitor models, they find a number of systematic trends (Fig. 8): the O shell typically has a higher convective Mach number (0.1-0.3) than the Si shell, where usually Ma \ 0 : 1, but there are islands around 16 M /C12 and 19 M /C12 in ZAMS mass where the convective Mach number in the Si shell reaches about 0.15 and is higher than in the O shell. The highest convective Mach numbers of up to 0.3 are reached in the O shell of lowmass progenitors with small cores as O burns deeper in the gravitational potential at higher temperatures. The general trend towards lower convective velocities in the O shell with higher progenitor and core mass is modified by variations in shell entropy and the residual O mass fraction at the onset of collapse. Deviations from this general trend also come about because the various C, Ne, O, and Si shell burning episodes do not always occur in the same order, and because of shell mergers.

The O shell is usually thicker and therefore allows large-scale modes with wave numbers ' \ 10 to dominate. Large-scale modes are more prevalent in progenitors above 16 M /C12 with their more massive O shells. The first, thick Si shell is no longer active at collapse in most cases, and there is typically only a thin convective Si shell

<!-- image -->

Fig. 8 Convective Mach number (left) and dominant angular wave number (right) in the Si shell (black) and O shell (red) predicted from 1D single-star evolution models from the study of Collins et al. (2018). Image reproduced with permission, copyright by the authors

<!-- image -->

(if any) between the Fe core and O shell at collapse, which will dominated by smallscale motions.

Collins et al. (2018) also find a high prevalence of late shell O-Ne shell mergers among high-mass progenitors. In about 40% of their models between 16 M /C12 and 26 M /C12 such a merger was initiated within the last minutes of collapse.

Although some of these trends follow from robust structural features and trends in the progenitor evolution, these findings will need to be examined with different stellar evolution codes and may be modified in detail, especially when better prescriptions for convective boundary mixing on secular time scales become available.

## 3.3 Convective boundaries

Mixing by entrainment As one of the most conspicuous features in their first 2D models of O shell burning, Bazan and Arnett (1994) and Baza ´n and Arnett (1997) noted the mixing of considerable amounts of C from the overlying layer into the active burning region. Although mixing across convective boundaries (sometimes indistinctly called ''overshooting'') had already been a long-standing topic in stellar evolution by then, these results were noteworthy because Bazan and Arnett (1994); Baza ´n and Arnett (1997) found much stronger convective boundary mixing (CBM) than compatible with overshoot prescriptions in 1D stellar evolution models of massive stars. Second, they observed that the mixed material can burn vigorously and thereby in turn dramatically affect the convective flow, i.e., there is the possibility of a feedback mechanism that cannot occur in the case of envelope convection or surface convection. Meakin and Arnett (2006) investigated this problem further in a situation with active and interacting O and C shells and observed strong excitation of p- and g-modes at convective-radiative boundaries, which, as they suggested, might also contribute to compositional mixing.

Critical steps beyond a mere descriptive analysis of CBM during the late burning stages were finally taken by Meakin and Arnett (2007a), who established i) the presence of CBM also in 3D (albeit weaker than in 2D), ii) identified the dominant process as entrainment driven by shear (Kelvin-Helmholtz and Holmbo ¨e) instabilities at the convective boundary, and iii) verified that the mass entrainment \_ M rate obeys a power law that can be motivated theoretically and has been verified in laboratory experiments of shear-driven entrainment (Fernando 1991; Strang and Fernando 2001):

<!-- formula-not-decoded -->

Here, A and n are dimensionless constants and Rib is the bulk Richardson number, which can be expressed in terms of the integral scale L of the turbulent flow and the buoyancy jump D b across the boundary,

<!-- image -->

<!-- formula-not-decoded -->

The buoyancy jump can be obtained by integrating the square of the Brunt-Va ¨isala over the extent of the boundary layer from r 1 to r 2 ,

Z

<!-- formula-not-decoded -->

In the case of a thin boundary layer, this reduces to D b ¼ g dq = q , where dq = q is the density contrast across the convective interface. From their simulations, Meakin and Arnett (2007b) determined values of A ¼ 1 : 06 and n ¼ 1 : 05 for the power-law coefficients. Since the work expended to entrain material against buoyancy the force of buoyancy must be supplied by a fraction of the convective energy flux (an argument which was independently redeveloped by Spruit 2015), one expects a value of n ¼ 1 for sufficiently high RiB.

Several subsequent 3D simulations (Mu ¨ller et al. 2016b; Jones et al. 2017) have confirmed a value of n /C25 1 for the scaling law (38). Mu ¨ller et al. (2016b) found a significantly smaller value of A /C25 0 : 1, however, but this may simply be due to ambiguities in the definition and measurement of the integral length scale L , which Mu ¨ller et al. (2016b) identify with the pressure scale height K , and of the convective velocity v conv that enters Eq. (39) for the bulk Richardson number. Jones et al. (2017) expressed the entrainment law slightly differently by a proportionality \_ M / \_ Q nuc to the total nuclear energy generation rate \_ Q nuc , which is equivalent to Eq. (38) with n ¼ 1. Their simulations are particularly noteworthy because they employed sufficiently high resolution to establish the entrainment law up to very high Rib. Although they do not explicitly state values of Rib, one can estimate that their models reach up to Rib ¼ 700-1000.

The simulations of Cristini et al. (2017, 2019) are a notable exception as they find a significantly shallower power law with n ¼ 0 : 74. This different power-law slope has yet to be accounted for, but it is important to note that despite the shallower power law, Cristini et al. (2017, 2019) generally find lower entrainment rates than Meakin and Arnett (2007b) for the same value of Rib with a much smaller value of A ¼ 0 : 05. At Rib /C25 20, their entrainment rate is actually in very good agreement with Mu ¨ller et al. (2016b), and in the region of Rib ¼ 40-300 their data are consistent with a steeper power law of n /C25 1. Since Cristini et al. (2017, 2019) also explore a much broader range in bulk Richardson number than the aforementioned studies, one possible interpretation could be that i) the value of A was overestimated in Meakin and Arnett (2007b), and that ii) the low value of n ¼ 0 : 74 may be due to a flattening of the entrainment law below Rib ¼ 20-30 for some physical reasons, and perhaps a slight flattening at Rib [ 200 because of numerical resolution effects.

Shell mergers Convective boundary mixing can take on a dramatic form when the buoyancy jump between two shells is sufficiently small for the neighboring shells to

<!-- image -->

merge entirely within a few convective turnover times. Such shell mergers have long been known to occur in 1D stellar evolution models, in particular between O, Ne, and C shells (e.g., Sukhbold and Woosley 2014; Collins et al. 2018). This is because balanced power leads to very similar entropies in the O, Ne, and C shells, and hence small buoyancy jumps between the shells. When nuclear energy generation and neutrino cooling finally fall out of balance due to shell contraction, the entropy of the inner (O or Ne) shell frequently increases and overtakes the outer shell(s), so that such mergers are particularly prevalent shortly before collapse as pointed out by Collins et al. (2018), who estimated that 40% of stars between 16 M /C12 and 26 M /C12 collapse during an ongoing shell merger. Although such mergers occur in 1D models, they may occur more readily in 3D, and 3D simulations are also necessary to capture the composition inhomogeneities and nucleosynthesis during the dynamical merger phase.

Shell mergers have indeed been seen in several recent 3D simulations. Mu ¨ller (2016) pointed out the breakout of a thin O shell through an inert, non-convective O layer into the active Ne burning zone in a 12 : 5 M /C12 star in the last minute before collapse, which, however, did not lead to a complete shell merger. Moca ´k et al. (2018) found a merger between the O and Ne shell in a 23 M /C12 model, and noted that the runaway entrainment leads to a peculiar quasi-steady with two distinct burning zones for O (at the base) and Ne (further out) within the same convective shell. However, their simulation only covered five turnover times and showed the merger occurring during the initial transient phase. Yadav et al. (2020) simulated an O-Ne shell merger in an 18 : 88 M /C12 over 15 turnover times, and were able to follow the evolution from the pre-merger phase with a soft, but clearly defined shell boundary and slow steady-state entrainment through the dynamical merger phase to a partially mixed post-merger state at the onset of collapse.. They stressed the emergence of large-scale asymmetries in the velocity field (with extreme velocities of up to 1700 km s /C0 1 ) and the composition during the merger, although the compositional asymmetries are already washed out somewhat at the point of collapse.

Impact on nucleosynthesis With multi-D simulations of the late-burning stages firmly established, it is critical to identify observable fingerprints of additional convective boundary mixing. The nucleosynthesis yields may provide one such fingerprint, which has already been discussed by several studies, even though one can only draw conclusions based on qualitative arguments and on 1D models with artificially enhanced mixing so far.

Davis et al. (2019) pointed out that the assumptions for convective boundary mixing can significantly affect the yields of various a -elements (C, O, Ne, Mg, Si), simply as a consequence of the change in shell structure. However, entrainment and shell mergers may leave more specific abundance patterns. In their investigation of O-C shell mergers in 1D Ritter et al. (2018) found significant overproduction of P, Cl, K, Sc, and possibly p-process isotopes, and argue that the occurrence of shell mergers may have important consequences for galactic chemical evolution (GCE). More recently, Co ˆte ´ et al. (2020) considered a Si-C shell merger, for which they find significant overproduction of 51 V and 52 Cr, which allows them to strongly constrain the rate of such events based on observed Galactic stellar abundances.

<!-- image -->

This is related to the long-standing realization that the ashes of hydrostatic silicon burning under neutron-rich conditions cannot be ejected in large quantities because of GCE constraints (Woosley et al. 1973; Arnett 1996).

Supernova spectroscopy may also help constrain additional convective boundary mixing and shell mergers may via their nucleosynthetic fingerprints. For example, the ejection of neutron-rich material from the silicon shell that is mixed out by entrainment before the explosion would lead to a supersolar Ni/Fe ratio as observed in some supernovae (Jerkstrand et al. 2015). Mixing of minimal amounts of Ca into O-rich zones can also have significant repercussions since only a mass fraction of a few 10 /C0 3 in Ca is required for Ca to be the domnant coolant during the nebular phase and quench O line mission in a shell (Fransson and Chevalier 1989; Kozma and Fransson 1998). This diagnostic potential of supernova spectroscopy for convective boundary mixing needs to be explored further in the future, but further (macroscopic) mixing during the explosion presents a major complication as it is not straightforward to disentangle the effect of mixing processes prior and during the explosion.

Secular effect on stellar evolution Evaluating the observational consequences of the convective boundary mixing seen in 3D models is also difficult because there is still no rigorous method for treating these processes in 1D stellar evolution codes. A crude estimate for the shell growth by entrainment can be formulated by noting that the work required to entrain material with density contrast dq = q against buoyancy must be no larger than a fraction of the time-integrated convective energy flux (Spruit 2015; Mu ¨ller 2016). During the late burning stage, the convective energy flux is set by the nuclear energy generation rate, and hence one can estimate that the entrained mass D M entr over the lifetime of a shell with mass M shell and radius r is roughly

<!-- formula-not-decoded -->

where A is the dimensionless coefficient in Eq. (38) and D Q is the average Q -value for a given burning stage (Mu ¨ller 2016). Based on Eq. (41), Mu ¨ller (2016) estimates that O shells could grow by tens of percent in mass by entrainment; for Si shells one expects a smaller effect, for C shells, the effect may be bigger.

How one can go beyond such simple estimates by using improved recipes for convective boundary mixing in stellar evolution codes is still an unresolved question. A common approach, based on the simulations of surface convection by Freytag et al. (1996), models entrainment as diffusive overshooting with an exponential decay of the MLT diffusion coefficient outside the convective boundary. The length scale k OV for the exponential decay can be calibrated against 3D simulations. This approach has been followed by Co ˆte ´ et al. (2020), Davis et al. (2019) and by many works on additional convective boundary mixing in low-mass mass stars, but has several issues. Entrainment is a very different process than diffusive overshoot that operates in a distinct physical regime (high Pe ´clet number), and hence one should not expect that it can be described by the same formalism (Viallet et al. 2015). It is also unclear why diffusive overshoot should be applied

<!-- image -->

only for compositional mixing in the entrainment regime. The common approach of expressing k OV by the pressure scale height is also open to criticism because the relevant length scale should be set by the convective velocities and the buoyancy jump, so that one would rather expect k OV / v 2 conv = D b .

Staritsin (2013) proposed an alternative approach of extending convectively mixed regions with time following the entrainment law (38), which better reflects the physics of the entrainment process. However, this approach has not been applied to the late neutrino-cooled burning stages yet (i.e., precisely where entrainment should operate). It also has some conceptual issues, for example, the entrainment law (38) obviously breaks down if there is convection on both sides of a shell interface. Yet another approach for extra mixing in 1D models was followed by Young et al. (2005), who handle mixing based on the local gradient Richardson number Ri for shear flows, which is estimated using an elliptic equation for the amplitudes of waves excited by convective motions (Young and Arnett 2005). This approach is physically motivated, but is still awaiting (and worthy of) a more quantitative comparison with 3D simulations beyond the qualitative discussion in Young and Arnett (2005).

Flame propagation in low-mass progenitors Around the minimum progenitor mass for supernova explosions, multi-D effects can have a more profound effect than merely changing shell masses on a modest scale, and may decide about the final fate of the star. This regime is best exemplified by the electron-capture supernova channel of super-asymptotic giant branch (SAGB) stars (see Jones et al. 2013; Doherty et al. 2017; Nomoto and Leung 2017; Leung and Nomoto 2019 for a broader overview and a discussion of uncertainties). In this evolutionary channel, the star builds up a degenerate core composed primarily of O and Ne. If this core grows to 1 : 38 M /C12 , electron captures on Ne and Mg trigger an O deflagration. Depending on the interplay of deleptonization (which decreases the degeneracy pressure) and the nuclear energy release, the core either contracts, collapses to a neutron stars, and explodes as an electron capture supernova, or the core does not collapse and explodes as a weak thermonuclear supernova (Jones et al. 2016). Since the flame is turbulent, its propagation needs to be modelled in multi-D, similar to deflagrations in Type Ia supernovae. Simulations of this problem have been conducted by Jones et al. (2016), Kirsebom et al. (2019) in 3D and (Leung and Nomoto 2019; Leung et al. 2020) in 2D. Efforts to improve the nuclear input physics and explore the sensitivity to the ignition geometry, general relativity, and flame physics are ongoing (e.g., Kirsebom et al. 2019; Leung et al. 2020).

For slightly more massive stars, one encounters similar situations with convectively-bounded flames after off-center ignition of O or Si (Woosley and Heger 2015). Again, multi-D effects may significantly affect the final evolutionary phase before collapse in this regime, but multi-D simulations of such supernova progenitors are yet to be carried out (but see Lecoanet et al. 2016 for idealized 3D simulations relevant to this regime).

<!-- image -->

## 3.4 Current and future issues

Significant progress notwithstanding, multi-D simulations of the late stages of convective burning still face further challenges. For the evolution towards collapse, models will eventually need to include a more sophisticated treatment of burning and deleptonization in the QSE and NSE regime and forgo the current approach of either using small networks or excising the Fe/Si core. Perhaps an even greater concern about the conclusions of current multi-D simulations lies with the timescale problem, however. No 3D simulations have yet been evolved sufficiently long to reach the state of balanced power (or to reveal why it would not be reached). This may have repercussions for turbulent entrainment, which ultimately taps the energy in turbulent motions and hence cannot be completely disconnected from the energy budget within a shell.

Moreover, although a few attempts have been made to simulate convection in rotating shells in 2D (Arnett and Meakin 2010; Chatzopoulos et al. 2016) and 3D (Kuhlen et al. 2003), multi-D simulations have yet to investigate the angular momentum distribution and angular momentum transport during the late presupernova stages in a satisfactory manner. Three-dimensional simulations are even more critical for this purpose than in the non-rotating case since many relevant phenomena such as (Rossby waves, Taylor-Proudman columns) cannot be modeled adequately in 2D. Simulations also need to explore a larger parameter space since the convective dynamics will depend on the Rossby number Ro /C24 v conv = ð X R Þ (where X is the rotational velocity). Furthermore, the time scales become a bigger challenge because models need to be run for several rotational periods T ¼ 2 pX /C0 1 and several convective turnover times s conv (whichever is longer), which is problematic since rotation in pre-supernova models is likely slow (e.g., Heger et al. 2005) so that Ro /C29 1 and T /C29 s conv. On the bright side, multi-D simulations may reveal much more interesting differences to 1D stellar evolution models: Both Kuhlen et al. (2003) and Arnett and Meakin (2010) found pronounced differential rotation developing from a rigidly rotating initial state, and Arnett and Meakin (2010) suggest that convective shells might adjust to a stratification with constant angular momentum as a function of radius rather than to uniform rotation as assumed in stellar evolution models. However, more simulations and more rigorous analysis is required to investigate these claims.

The problem of rotation can obviously not be solved without including magnetic fields in the long run. It is well known (e.g., Shu 1992) that the criterion for the instability of rotating flow becomes less restrictive in the MHD case, and effects such as the magnetorotational instability (MRI, Balbus and Hawley 1991) or a small-scale dynamo may enforce a more uniform rotation profile than hydrodynamic convection alone. But the importance of magnetic fields is not confined to the case rotating progenitors. Prompted by helioseismic measurements that indicate smaller convective velocities in the deeper layer of the solar convection zone (Gizon and Birch 2012; Hanasoge et al. 2012), some simulations of magnetoconvection in the Sun found a suppression of convective velocities by up to 50% compared to hydrodynamic simulations due to strong magnetic fields from a small-scale dynamo

<!-- image -->

that reach equipartition with the turbulent kinetic energy (Hotta et al. 2015). Magnetic fields can also inhibit or enhance mixing in shear layers (Bru ¨ggen and Hillebrandt 2001) and may hence affect convective boundary mixing. Thus, there is still plenty of ground left to explore for simulations of the late burning stages.

## 4 Core collapse and shock revival

In the Introduction, we already outlined a variety of multi-D effects than can play a role in reviving the stalled supernova shock as a subsidiary agent to neutrino heating (i.e., neutrino-driven convection in the gain region, the SASI, and progenitor asphericities), or also as the main drivers of the explosion (rotation and magnetic fields). Historically, a number of works have also considered convection in the PNS interior as a means for precipitating explosions by enhancing the neutrino emission from the PNS (Epstein 1979; Wilson and Mayle 1988; Burrows and Lattimer 1988; Wilson and Mayle 1993), but these hopes were not substantiated in subsequent decades. Nonetheless, convection inside the PNS remains important for various aspects of the supernova problem such as the neutrino and gravitational wave signals and the nucleosynthesis conditions in the innermost ejecta.

Since each of these phenomena has proved rich and complex over the years, it is no longer possible to treat them adequately within a chronological narrative of the quest for the supernova explosion mechanism. Nevertheless, ascertaining the explosion mechanism by means of first-principle simulations remains the overriding concern in supernova theory, and it is therefore still useful to recapitulate the progress in supernova explosion modelling from the advent of the first 2D models with a simplified treatment of neutrino heating and cooling in the 1990s (Herant et al. 1992; Shimizu et al. 1993; Yamada et al. 1993; Janka et al. 1993; Herant et al. 1994; Burrows et al. 1995; Janka and Mu ¨ller 1995, 1996). A more detailed analysis of the individual hydrodynamic phenomena beyond this chronicle of simulations is then provided in Sects. 4.1-4.7.

Neutrino-driven explosions in 2D Although the 2D simulations of the early and mid1990s had shown multi-D effects to be helpful for shock revival, these models did not utilize neutrino transport on par with the best available methods for 1D simulations at the time. In a first attempt to better model neutrino heating and cooling in 2D by using the pre-computed neutrino radiation field from a 1D simulation with multi-group flux limited diffusion, Mezzacappa et al. (1998) were unable to reproduce the successful explosions found in earlier 2D models. This led to a resurgence of interest in accurate methods for neutrino transport, culiminating in the development of Boltzmann solvers for relativistic (Yamada et al. 1999; Liebendo ¨rfer et al. 2001, 2004) and pseudo-Newtonian simulations (Rampp and Janka 2000, 2002). The explosion problem was then revisited in 2D using various types of multi-group neutrino transport from the mid-2000s onwards. Neutrinodriven explosions were obtained in many of these 2D simulations for a wide range of progenitors (Buras et al. 2006a; Marek and Janka 2009; Mu ¨ller et al. 2012b, a, 2013; Janka 2012; Janka et al. 2012; Suwa et al. 2010, 2013; Bruenn

<!-- image -->

et al. 2013, 2016; Nakamura et al. 2015; Burrows et al. 2018; Pan et al. 2018; O'Connor and Couch 2018b), though still with significant differences between the various simulation codes.

Challenges and successes in 3D Following isolated earlier attempts at 3D modelling using the ''light-bulb'' style models of the 1990s (Shimizu et al. 1993; Janka and Mu ¨ller 1996), or gray flux-limited diffusion (Fryer and Warren 2002), the role of 3D effects in the explosion mechanism was finally investigated vigorously in the last decade, starting again with simple light-bulb models (Nordhaus et al. 2010b; Hanke et al. 2012; Couch 2013; Dolence et al. 2013). Except for spurious results in Nordhaus et al. (2010b), these light-bulb models indicated a similar ''explodability'' in 2D and 3D. However, subsequent 3D models with rigorous neutrino transport proved more reluctant to explode; indeed the first 3D models of 11 : 2 M /C12 , 20 M /C12 , and 27 M /C12 progenitors using multi-group, three-flavour neutrino transport did not explode at all (Hanke et al. 2013; Tamborra et al. 2014a).

Even though various groups have now obtained explosions in 3D simulations, shock revival usually occurs later than in 2D, and often requires additional (and sometimes hypothetical) ingredients to improve the heating conditions or a specific progenitor structure. For low-mass single-star (Melson et al. 2015b; Mu ¨ller 2016; Burrows et al. 2019) and binary (Mu ¨ller et al. 2018) progenitors just above the ironcore formation limit, 3D simulations readily yield explosions since the steep drop of the density outside the iron core implies a rapid drop of the accretion rate onto the shock after bounce. For more massive stars the record is mixed. For standard, nonrotating progenitors in the range between 11 : 2 M /C12 and 27 M /C12 and unmodified, stateof-the-art microphysics, no explosions were found in simulations using the VERTEX code (Hanke et al. 2013; Tamborra et al. 2014a; Melson et al. 2015a; Summa et al. 2018) and the FLASH-M1 code (O'Connor and Couch 2018a). On the other hand, the Oak Ridge group obtained an explosion for a 15 M /C12 star (Lentz et al. 2015) with their CHIMERA code, and the Princeton group observed shock revival in eleven out of fourteen models between 9 M /C12 and 60 M /C12 (Vartanyan et al. 2019b; Burrows et al. 2019; Radice et al. 2019; Burrows et al. 2020) with the FORNAX code. In both cases the accuracy of the microphysics, the neutrino transport, and gravity treatment appears comparable to VERTEX. Three-dimensional simulations using other codes (that are constantly evolving!) are more difficult to compare as they involve simplifications in the microphysics or transport compared to VERTEX, CHIMERA, and FORNAX, although some of them compensate for this by higher resolution in real space and energy space and a better treatment of gravity. At any rate, results obtained with other codes such as COCONUT-FMT, FUGRA, ZELMANI, and 3DNSNE add to the picture of simulations straddling the verge between successful shock revival (Takiwaki et al. 2012, 2014; Mu ¨ller 2015; Roberts et al. 2016; Chan et al. 2018; Ott et al. 2018; Kuroda et al. 2018) and failure (Mu ¨ller et al. 2017a; Kuroda et al. 2018) for standard, non-rotating progenitors and standard or simplified microphysics.

These different results may simply indicate that the neutrino-driven mechanism operates close to the critical threshold for explosion in nature. Observations of supernova progenitors indeed indicate that black hole formation occurs already at relatively low masses down to to /C24 15 M /C12 for single stars (Smartt et al. 2009;

<!-- image -->

Smartt 2015). Since the lack of robust explosions in 3D persists to even lower masses, and since strongly delayed explosions in 3D may turn out too weak to be compatible with observations, several groups have explored new avenues towards more robust explosions. Some of the proposed ideas invoke modifications or improvements to the microphysics that ultimately lead to improved neutrino heating conditions, such as strangeness corrections to the neutral-current scattering rate (Melson et al. 2015a) and muonization (Bollig et al. 2017). Other studies have explored purely hydrodynamic effects. Among these, Takiwaki et al. (2016), Janka et al. (2016), Summa et al. (2018) pointed out that rapid progenitor rotation could be conducive to shock revival even without invoking MHD effects. Another idea posits that including seed perturbations from the late convective burning stages can facilitate shock revival. First studied by Couch and Ott (2013) and Mu ¨ller and Janka (2015) using parametric initial conditions, this ''perturbation-aided'' mechanism has subsequently been explored further using pre-collapse perturbations from 3D models of the late burning stages, initially with ambiguous results in the leakage simulations of Couch et al. (2015) and then in a number of 3D simulations using multi-group neutrino transport (Mu ¨ller 2016; Mu ¨ller et al. 2017a, 2019), where it led to robust explosions over a wider mass range from 11 : 8 M /C12 to 18 M /C12 for single stars.

Neutrino-driven explosion models have thus matured considerably in recent years, but it would be premature to declare the problem of shock revival solved. The discrepancies between the results of different groups have yet to be sorted out, and there is still no ''gold standard'' among the simulations that combines the best neutrino transport, the best microphysics, 3D progenitor models, and general relativity. Moreover, phenomenological models of neutrino-driven explosions (Ugliano et al. 2012; Pejcha and Thompson 2015; Sukhbold et al. 2016; Mu ¨ller et al. 2016a) suggest that a different mechanism is still needed to explain hypernova explosions with energies above 2 /C2 10 51 erg.

Magnetohydrodynamic simulations The mechanism(s) behind hypernovae likely rely on rapid rotation and strong magnetic fields (Akiyama et al. 2003; Woosley and Bloom 2006), but the importance of magnetic fields may not end there. There may be a continuous transition from neutrino-driven explosions to MHD-driven explosions (Burrows et al. 2007a), and strong magnetic fields may also a role in non-rotating progenitors as an important driving agent or as a subsidiary to neutrino heating (Obergaulinger et al. 2014).

Although the ideas of Akiyama et al. (2003) quickly triggered first 2D MHD core-collapse supernova simulations (e.g., Yamada and Sawai 2004; Sawai et al. 2005; Obergaulinger et al. 2006; Shibata et al. 2006), there is still only a small corpus of magnetorotational supernova explosion models, especially if we focus on models of the entire collapse, accretion, and early explosion phase using reasonably detailed microphysics and disregard parameterized models of relativistic and nonrelativistic jets and of collapsar disks. Burrows et al. (2007a) presented 2D simulations of magnetorotational explosions of a 15 M /C12 progenitor (later followed by MHD simulations of accretion-induced by collapse in Dessart et al. 2007) with the Newtonian radiation-MHD code VULCAN and demonstrated the ready emergence

<!-- image -->

of jets powered by strong hoop stresses for sufficiently strong initial fields. Burrows et al. (2007a) made the important point that these non-relativistic jets are a distinctly different phenomenon from the relativistic jets seen in long gamma-ray bursts (GRBs), which may be formed several seconds after shock revival.

Like most other subsequent simulations, these models relied on parameterized initial conditions with artificially strong magnetic fields to mimic the purported fast amplification of much weaker fields in the progenitor by the magnetorotational instability (Balbus and Hawley 1991; Akiyama et al. 2003). They also imposed the progenitor rotation profile by hand. The 2D studies of Obergaulinger and Aloy (2017, 2020), Bugli et al. (2020) have recently explored variations in the assumed initial field strength and topology and the assumed rotation profiles more thoroughly. While they find considerable variation in the outcome of their models, it is interesting to note that in some instances Obergaulinger and Aloy (2020) even find magnetorotational explosions for the unmodified rotation profile and magnetic field strength of two of the 35 M /C12 progenitor models from Woosley and Heger (2006), although it is not perfectly clear what the precise geometry of the field in the stellar evolution models ought to be.

The imposition of axisymmetry is an even bigger concern in the case of magnetorotational supernovae than for nuetrino-driven explosion models. Several 3D simulations of MHD-driven explosions have by now been performed, but among these only Obergaulinger and Aloy (2020) included multi-group neutrino transport, whereas the others (Winteler et al. 2012; Mo ¨sta et al. 2014b, 2018) employed a leakage scheme. The prospects for successful magnetorotational explosions in 3D are still somewhat unclear. Mo ¨sta et al. (2018) reported the destruction of the emerging jets by a kink instability, although the jet can apparently be stabilised if the poloidal field strength is comparable to the toroidal field strength (Mo ¨sta et al. 2018). Moreover, the explosion dynamics already depends sensitively on the assumed initial field geometry; strong dipole fields appear to be required for the most powerful explosions (Bugli et al. 2020). Given the vast uncertainties concerning the initial rotation rates, field strengths, and field geometries in the supernova progenitors, considerably more work is necessary before the magnetorotational mechanism can be considered robust even for a small sub-class of progenitors. We will therefore focus only on the hydrodynamics of neutrino-driven explosions in the subsequent discussion.

## 4.1 Structure of the accretion flow and runaway conditions in spherical symmetry

Before analyzing the role of multi-D phenomena in core-collapse supernovae in greater depth, it is expedient to discuss the structure of the supernova core that emerges once the gain region has formed a few tens of milliseconds after collapse in an idealized, spherically-symmetric picture as shown in Fig. 9. Our discussion closely follows the works of Janka (2001), Mu ¨ller and Janka (2015) and Mu ¨ller et al. (2016a) which may be consulted for further details.

<!-- image -->

Fig. 9 Schematic 1D structure of the supernova core after the formation of the gain region, illustrated by profiles of the density q , pressure P , temperature T (left), radial velocity vr , entropy s , and electron fraction Y e (right). The profiles are taken from a 1D radiation hydrodynamics simulation of the 20 M /C12 progenitor of Woosley and Heger (2007) at a post-bounce time of 200 ms. See text for details

<!-- image -->

Structure of the accretion flow At this stage, the PNS consists of an inner core of about 0 : 5 M /C12 (depending on the equation of state) with low entropy, which is surrounded by an extended mantle of about 1 M /C12 that was heated to entropies of about 6 k b = nucleon as the shock propagated through the outer part of the collapsed iron core and most of the Si shell. The mantle extends out to the neutrinosphere at high subnuclear densities, where neutrinos on average undergo their last interaction before escape (for details see Kotake et al. 2006; Janka 2017; Mu ¨ller 2019b; Mezzacappa 2020). In the atmosphere immediately outside the neutrinosphere radius R m , the pressure P is dominated by non-relativistic baryons, and neutrino interactions are still frequent enough to act as ''thermostat'' and maintain a roughly isothermal stratification, resulting in an exponential density profile (Janka 2001):

/C20

/C18

/C19

/C21

<!-- formula-not-decoded -->

where M is the PNS mass, R m , T m , and qm are the neutrinosphere radius, temperature, and density, and m n is the neutron mass. To maintain rough isothermality with the neutrinosphere, the accreted matter must cool as it is advected through the atmosphere. Below a density of about 10 10 g cm /C0 3 , the pressure is dominated by relativistic electron-positron pairs and photons, and around this point neutrino heating starts do dominate over neutrino cooling at the gain radius R g . 7 Since the cooling and heating rate scale with T 6 / P 3 = 2 and L m E 2 m = r 2 in terms of the matter temperature T and the electron-flavor neutrino luminosity L m and mean energy E m (appropriately averaged over electron neutrinos and antineutrinos), balance between heating and cooling defines an effective thermal boundary condition for the radiation-dominated gain region further out,

7 Properly speaking, the EoS transition radius between the baryon-dominated and the radiationdominated regime and the gain radius are close, but the gain radius is slightly larger (Janka 2001). For many purposes it is not critical to distinguish them.

<!-- image -->

<!-- formula-not-decoded -->

Before shock revival, the stratification between the gain radius is roughly adiabatic out to the shock, 8 resulting in power-law profiles q / r /C0 3 and T / r /C0 1 for the temperature and density. Ahead of the shock, the infalling material moves with a radial velocity of j vr j /C25 ffiffiffiffiffiffiffiffiffiffiffiffi GM = r p (i.e., a large fraction of the free-fall velocity), and the density is given in terms of the mass acrretion rate \_ M as q ¼ \_ M = ð 4 p r 2 j vr jÞ . In a quasi-stationary situation, the stalled accretion shock will adjust to a radius R sh such that the jump conditions are fulfilled and the post-shock pressure P sh and the preshock ram pressure P ram ¼ q v 2 r are related by

r

ffiffiffiffiffiffiffi ffi

<!-- formula-not-decoded -->

Here b is the compression ratio in the shock, which varies from b /C25 10 early on, which is slightly larger than the value of b ¼ 7 for an ideal gas with adiabatic index c ¼ 4 = 3 because of the dissociation of nuclei in the shock, to b /C25 4 during the explosion phase when there is a net release of energy by burning in the shock. Equation (44) immediately implies that the quasi-stationary accretion shock radius increases with the post-shock pressure roughly as R sh / P 2 = 5 sh . Recognizing that the heating rate \_ q heat / L m T 2 m r /C0 2 and the cooling rate \_ q cool / T 6 / P 3 = 2 balance at the gain radius, and using the adiabatic stratification in the gain region and the jump conditions, one can go further and derive that the shock radius scales as (Janka 2012; Mu ¨ller and Janka 2015)

<!-- formula-not-decoded -->

in spherical symmetry in terms of L m , T m , R g , M and the mass accretion rate \_ M , which is related to the density profile of the progenitor (Woosley and Heger 2012; Mu ¨ller et al. 2016a).

Conditions for shock revival So far, we have assumed a stationary accretion flow in this picture. The problem of shock revival is, however, related to the breakdown of stationary accretion solutions (Burrows and Goshy 1993; Janka 2001), or more strictly speaking, to the development of non-linearly unstable flow perturbations (Ferna ´ndez 2012). The transition to runaway shock expansion can be understood in terms of a competition of time scales, namely the advection or residence time s adv that the accreted material spends in the gain region, and the time scale s heat for unbinding the material in the gain region by neutrino heating (Thompson 2000;

8 This is because neutrino heating does not change the entropy appreciably as material traverses the gain region as long as the heating conditions are far from critical. Furthermore, mixing reduces the entropy gradient in 3D once convection or SASI have developed.

<!-- image -->

Janka 2001; Thompson et al. 2005; Buras et al. 2006a; Murphy and Burrows 2008). These can be computed in terms of the binding energy E g and mass M g of the gain region, the mass accretion rate \_ M , and the volume-integrated heating rate \_ Q m as

<!-- formula-not-decoded -->

<!-- formula-not-decoded -->

Transition to runaway expansion is expected if s adv = s heat J 1, which is borne out by 1D light-bulb simulations with a pre-defined neutrino luminosity (Ferna ´ndez 2012). Alternatively, the runaway condition can be expressed in terms of a critical luminosity L crit above which there are no stationary 1D accretion solutions (Burrows and Goshy 1993). Janka (2012) and Mu ¨ller and Janka (2015) have pointed out that these two descriptions are essentially equivalent by converting the time-scale criterion into a power law for the critical value of the ''heating functional'' L ¼ L m E 2 m ,

<!-- formula-not-decoded -->

Other largely equivalent ways to characterize the onset of a runaway instability (at least in spherical symmetry) are the notion that the Bernoulli parameter reaches zero somehwere in the gain region around shock revival (Burrows et al. 1995; Ferna ´ndez 2012), and the antesonic condition c 2 s [ 3 = 8 GM = r (Pejcha and Thompson 2012), which is effectively a condition for the flow enthalpy just like the Bernoulli parameter (Mu ¨ller 2016).

## 4.2 Impact of multi-dimensional effects on shock revival

Qualitative description This simple picture is useful for qualitatively understanding how multi-D effects modify the spherically-averaged bulk structure of the postshock flow and hence affect the conditions for shock revival. It is intuitive from Eq. (44) that increasing the post-shock pressure (e.g., by turbulent heat transfer), or adding turbulent or magnetic stresses will increase the shock radius and modify Eq. (45) for the spherically symmetric case. This will then affect the time scales s adv and s heat and thereby modify the conditions for runaway shock expansion driven by neutrinos. Moreover, certain multi-D phenomena may also facilitate runaway shock expansion more directly by dumping extra energy into the gain region, which may take the form of thermal energy, turbulent kinetic energy, or magnetic energy. This is, of course, only a coarse-grain interpretation of the effect of multi-D effects, which needs to be based on a more careful analysis of the underlying hydrodynamic phenomena.

<!-- image -->

The studies from the 1990s also outlined qualitative explanations for the beneficial role of multi-D effects. Herant et al. (1994) interpreted convection as of an open-cycle heat engine that continuously pumps transfers energy from the gain radius (where neutrino heating is strongest) further out into the gain region, and Janka and Mu ¨ller (1996) similarly stress the importance of more effective heat transfer from the gain region to the shock. Herant et al. (1994) argued that largescale mixing motions are also advantageous because they continue to channel fresh matter to the cooling region during the explosion phase so that the neutrino heating is not quenched when the shock is revived. Finally, Burrows et al. (1995) pointed out what we now subsume under the notion of turbulent stresses: As the convective bubbles collide with the shock surface with significant velocities (or in modern parlance provide ''turbulent stresses'') and thereby deform and expand it.

Modified critical luminosity Since then, the impact of multi-D effects has been analyzed more quantitatively. Several studies (Buras et al. 2006a; Murphy and Burrows 2008; Hanke et al. 2012) showed that the advection time scale s adv is systematically larger in multi-D, while the onset of runaway shock expansion is still determined by the criterion s adv = s heat . This suggests that the runaway is still powered by neutrino heating just as in 1D, and that multi-D effects facilitate explosions facilitate shock revival by somewhat expanding the stationary shock, keeping a larger amount of mass in the gain region, and thereby increasing the heating efficiency. 9 To a lesser extent, mixing also reduces the binding energy of the gain region (Mu ¨ller 2016), but this appears to be of secondary importance for shock revival.

Building on the 1D picture from Sect. 4.1, the increase of the quasi-stationary shock radius can be understood as the consequence of additional ''turbulence'' 10 terms that arise in a spherical Reynolds or Favre decomposition of the flow, whose importance can be gauged by the square of the turbulent Mach number Ma in the gain region (Mu ¨ller et al. 2012a). Using light-bulb simulations, Murphy et al. (2013) demonstrated quantitatively that the inclusion of Reynolds stresses (' 'turbulent pressure'') largely accounts for the higher shock radius in multi-D models. The critical role of the turbulent pressure was confirmed by Couch and Ott (2015) using a leakage scheme and by Mu ¨ller and Janka (2015) with multi-group neutrino transport.

The resulting effect on the critical luminosity can be estimated by including a turbulent pressure term 11 in Eq. (44), which ultimately leads to (Mu ¨ller and Janka 2015)

9 Note that this refers to a comparison of multi-D and 1D models for a given set of parameters of the accretion flow ( L m , E m , M , \_ M , and R g. When comparing multi-D and 1D models at the threshold to explosion (with different L m and E m ), the heating efficiency can be lower in multi-D (Couch and Ott 2015), but this does not mean that there is a different runaway mechanism (Mu ¨ller 2016).

10 It is important to stress that ''turbulence'' is something of a convenient misnomer in this context and refers to any deviation from quasi-stationary, spherically-symmetric flow. This should be carefully distinguished from the usual notion of turbulence in high-Reynolds number flow, although the two concepts are frequently conflated.

11 An alternative approach to account for the effect of the turbulent pressure is to use an effective adiabatic index c [ 4 = 3 in the gain region (Radice et al. 2015).

<!-- image -->

/C18

/C19

/C18

/C19

<!-- formula-not-decoded -->

where the critical luminosity in 1D, ð L m E 2 m Þ crit ; 1D , is modified by a correction factor containing the turbulent Mach number in the gain region. Although based on a rather simple analytic model, Eq. (49) describes shock revival in 2D (Summa et al. 2016) and 3D models (Janka et al. 2016) remarkably well. This suggests that the critical parameter for increased explodability in multi-D is indeed the turbulent Mach number, although, as argued by Mabanta and Murphy (2018), the larger accretion shock radius may not be due to turbulent pressure alone. Even if other effects such as turbulent heat transport, turbulent dissipation (Mabanta and Murphy 2018), and even turbulent viscosity (Mu ¨ller 2019a) play a role, one expects a scaling law similar to Eq. (49) simply because any leading-order correction to the 1D jump condition (44) from a spherical Reynolds decomposition will scale with Ma 2 , only with a slightly different proportionality constant than in Eq. (49). The turbulent Mach number itself will be determined by the growth and saturation mechanisms of the non-radial instabilities in the gain region as discussed in the following sections.

## 4.3 Neutrino-driven convection in the gain region

Convection in the gain region develops because neutrino heating establishes a negative entropy gradient. In many respects, this ' 'hot-bubble convection'' resembles convection on top of a quasi-hydrostatic spherical background structure as familiar from the earlier phases of stellar evolution, but there are subtle differences because the instability occurs in an accrretion flow.

Condition for instability Under hydrostatic conditions, the Ledoux criterion for convective instability can be written as (Buras et al. 2006b),

/C18

/C19

/C18

/C19

<!-- formula-not-decoded -->

in terms of the gradients of density, pressure, entropy s and electron fraction Y e . Using a local stability analysis for a displaced blob, one finds a growth rate Im x BV, where the Brunt-Va ¨isa ¨la ¨ frequency x BV is defined as

<!-- formula-not-decoded -->

In a stationary accretion flow, the radial derivatives can be expressed in terms of the time derivatives \_ s and \_ Y e and the advection velocity vr ,

<!-- image -->

"

/C18

/C19

/C18

/C19

<!-- formula-not-decoded -->

Once the gain region forms around 80-100 ms after bounce, the electron fraction gradient plays a minor role for stability, and the material in the gain region can be well described as a radiation-dominated gas with P / ð q s Þ 4 = 3 so that

<!-- formula-not-decoded -->

This has an important consequence: different from a hydrostatic background, the stability of a heated accretion flow (or outflow) depends on the sign of the advection velocity rather than on the profile of the heating function. Using the aforementioned scaling for the heating rate and assuming a linear velocity profile behind the shock, we obtain an estimate

/C18

/C19

/C18

/C19

<!-- formula-not-decoded -->

More importantly, however, advection can stabilize the flow against convection because perturbations only have a finite time to grow as they cross the gain region as pointed out by Foglizzo et al. (2006), who demonstrated that instability is regulated by a parameter v ,

Z

<!-- formula-not-decoded -->

Instability should only occur for v J 3. This has indeed been confirmed in a number of parameterized (Scheck et al. 2008; Ferna ´ndez et al. 2014; Ferna ´ndez 2015; Couch and O'Connor 2014) and self-consistent simulations (Mu ¨ller et al. 2012a; Hanke et al. 2013) in 2D and 3D.

Dominant eddy scale Similar to the situation in convective shell burning, the length scale of the most unstable linear mode is determined by the width of the gain region according to Eq. (37) (Foglizzo et al. 2006). In 3D this remains the characteristic length scale of convective eddies during the non-linear saturation stage (e.g., Hanke et al. 2013). Since the ratio of the shock and gain radius typically lies in the range R sh = R g ¼ 1 : 5-2 before the heating conditions become close to critical, convection is characterized by medium-scale eddies with angular wavenumbers ' /C25 4-8 during the accretion phase (Hanke et al. 2013; Couch and O'Connor 2014). Around and after shock revival, large-scale modes with ' ¼ 1 and ' ¼ 2 emerge. By contrast, 2D simulations of hot-bubble convection tend to develop large-scale ( ' ¼ 1 and ' ¼ 2) vortices during the non-linear stage (Hanke et al. 2012, 2013; Couch 2013; Couch

<!-- image -->

#

and O'Connor 2014) as a result of the inverse turbulent cascade of 2D turbulence (Kraichnan 1967).

Non-linear saturation The evolution towards shock revival typically proceeds over sufficiently long time scales for hot-bubble convection to reach a quasi-stationary state. Using 2D light-bulb simulations, Murphy et al. (2013) first demonstrated that this quasi-stationary state closely mirrors the situation in stellar convection (cf. Sect. 3.1), i.e., neutrino heating, buoyant driving, and turbulent dissipation balance each other (see also Murphy and Meakin 2011), and as a result the convective luminosity scales with the neutrino heating rate. Alternatively, the quasi-stationary state can be characterized by the notion of marginal stability; the flow adjusts itself that such that the v -parameter for the spherically averaged flow converges to h v i /C25 3 (Ferna ´ndez et al. 2014). Mu ¨ller and Janka (2015) showed that these properties of the non-linear stage result in a scaling law for the convective velocity that is completely analogous to Eq. (36). In 2D simulations, the velocity perturbations d v scale as

<!-- formula-not-decoded -->

in terms of the average mass-specific neutrino heating rate \_ q m in the gain region. In 3D, the convective velocities are slightly smaller (Mu ¨ller 2016),

<!-- formula-not-decoded -->

The smaller proportionality constant in 3D can be motivated by the tendency of the forward turbulent cascade to create smaller structures, which decreases the dissipation length and increases the dissipation rate of the flow.

Quantitative effect on shock revival Based on these scaling laws for the convective velocity, Mu ¨ller and Janka (2015) determined that convective motions should reach a characteristic squared turbulent Mach number Ma /C24 0 : 3-0.45 around the time of shock revival. Using this value in Eq. (49) for the modified critical luminosity, they predict a reduction of the critical luminosity by 15-25% due to convection, which is in the ballpark of the numerical results (Murphy and Burrows 2008; Hanke et al. 2012; Couch 2013; Dolence et al. 2013; Ferna ´ndez et al. 2014; Ferna ´ndez 2015)

One might also be tempted to use Eqs. (56) and (57) to explain the lower explodability of self-consistent 3D models compared to their 2D counterparts. The nature of the differences between 2D and 3D is more complicated, however, since the critical luminosity for shock revival is roughly equal in 2D and 3D light-bulb simulations. Evidently, there are effects that partly compensate for the smaller convective velocities in 3D in some situations: the forward turbulent cascade (Melson et al. 2015b) and the different behavior of the Kelvin-Helmholtz instability in 3D (Mu ¨ller 2015) affect the interaction between updrafts and downdrafts and can result in reduced cooling in 3D (Melson et al. 2015b). Moreover, compatible with earlier studies of the Rayleigh-Taylor instability for planar geometry (Yabe et al.

<!-- image -->

1991; Hecht et al. 1995), Kazeroni et al. (2018) found a faster growth of convective plumes and more efficient mixing in a planer toy model of neutrino-driven convection. Along similar lines, Handy et al. (2014) appealed to the higher volumeto-surface ratio of convective plumes in 3D to explain the reduced critical luminosity in 3D in their light-bulb simulations. It is plausible that these factors establish a similar critical luminosity threshold for shock revival in light-bulb simulations, but they do not explain the much more decisive effect of dimensionality in self-consistent simulations. One possible explanation lies in the fact that explosions in self-consistent models usually occur in a short non-stationary phase with a rapidly decreasing mass accretion rate and neutrino luminosity around the infall of the Si/O shell interface; under these conditions the more sluggish emergence of large-scale modes in 3D due to the forward cascade may delay or inhibit shock revival (Lentz et al. 2015). Moreover, the more rapid response of the mass accretion rate to shock expansion in 3D (Melson et al. 2015b; Mu ¨ller 2015) might be hurtful around shock revival because this effect can reduce the accretion luminosity and hence undercut neutrino heating before a runaway situation can develop.

Resolution dependence and turbulence Because of the turbulent nature of neutrinodriven convection, the spectral properties of the flow and the the resolution dependence in simulations have received considerable attention in the literature. Most self-consistent models with multi-group neutrino transport can only afford a limited resolution (about 1 : 5 /C14 -2 /C14 in angle and about 100 zones or less in the gain region) and do not reach a fully developed turbulent state, with Handy et al. (2014) going so far as to speak of ' 'perturbed laminar flow'' instead. Various authors (Abdikamalov et al. 2015; Radice et al. 2015, 2016) have argued that considerably

<!-- image -->

Fig. 10 The advective-acoustic mechanism for the standing accretion shock instability. Upward propagating acoustic waves (blue) generate vorticity perturbations (red) as they interact with the accretion shock (orange circle). The vorticity perturbations are advected downward with the accretion flow to the PNS surface where they generate acoustic waves due to advective-acoustic coupling in the steep density gradient. Instability for a given mode obtains if the product of the amplitude ratios Q sh and Q r for between ingoing and outgoing waves at the shock and PNS surface satisfies Q sh Q r [ 1. Image reproduced with permission from Guilet and Foglizzo (2012), copyright by the authors

<!-- image -->

higher resolution is needed to obtain clean turbulence spectra with a developed inertial range and a Kolmogorov spectrum and raised concerns that a pile-up of kinetic energy (' 'bottleneck effect'') at small scales might affect the overall dynamics. However, the detailed spectral properties of the flow are usually not critical, and integral properties of the flow are more important for the impact of convection on shock revival. The resolution dependence nonetheless remains a concern for the question of shock revival, as some 3D resolution studies (Hanke et al. 2012; Abdikamalov et al. 2015; Roberts et al. 2016) found a trend towards decreasing explodability with increased resolution. Recent work by Melson et al. (2020) resolved most of these concerns in a resolution study using light-bulb simulations. They demonstrated that the resolution dependence in Hanke et al. (2012) was a spurious effect connected to details in their light-bulb scheme, and instead found a trend towards increased explodability at higher resolution. In rough agreement with Handy et al. (2014), Melson et al. (2020) found that the overall flow dynamics converges at an angular resolution of about 1 /C14 , which is not far from what most self-consistent simulations can afford (but one still needs to bear in mind that the resolution requirements depend on the details of the numerical scheme, cf. Sect. 2.1.2). Melson et al. (2020) also pointed out that neutrino drag plays a non-negligible role in the gain region, so that merely increasing the resolution does not add physical realism beyond numerical Reynolds numbers of a few hundred unless neutrino drag is also included as a non-ideal effect. Melson et al. (2020) speculate that findings of decreased explodability with higher resolution in Cartesian 3D models may be explained because the grid-induced seed asphericities are lower at higher resolution.

## 4.4 The standing accretion shock instability

Using adiabatic 2D simulations of spherical accretion shocks, the seminal work of Blondin et al. (2003) demonstrated that another instability, dubbed ' 'SASI'' (standing accretion shock instability), can operate in the supernova core even without a convectively unstable gradient in the gain region. This instability takes the form of large-scale ( ' ¼ 1 and sometimes ' ¼ 2) oscillatory motions of the shock, and it was immediately realized that it can support shock revival in a similar manner as convection. In early 2D supernova simulations, the SASI was sometimes confused with convection because the two phenomena share superficial similarities like high-entropy bubbles and low-entropy accretion downflows. However, the SASI is set apart from convection by dipolar (and sometimes quadrupolar) flow. Foglizzo et al. (2006) pointed out that for the typical ratio between the shock and gain radius in the pre-explosion phase there are no unstable convective modes with ' ¼ 1 ; 2 in the gain region; instead one finds ' /C25 4-8 according to Eq. (37), and for v \ 3 the flow becomes stable against convection altogether without strong perturbations (see also Sect. 4.3). This implies that a different instability mechanism-the one discovered by Blondin et al. (2003)-must be responsible for the ' ¼ 1 and ' ¼ 2 modes in supernova models with small ratios R sh = R g .

Amplification mechanism The stability of accretion shocks had in fact already been analyzed earlier in the context of accretion onto compact objects using linear

<!-- image -->

perturbation theory (Houck and Chevalier 1992; Foglizzo 2001, 2002), which provided useful groundwork for identifying the physical mechanism behind the SASI and explaining the ' ¼ 1 ; 2 nature of the instability from its dispersion relation. The accepted picture is now that of a vortical-acoustic cycle (Foglizzo 2002; Foglizzo et al. 2007): Shock deformation generates vorticity waves that are advected towards the PNS surface, and due the deceleration of the flow in the steep density gradient below the gain region (see Fig. 9), these vorticity waves in turn generate outgoing sound waves that again couple to acoustic waves at the shock (Fig. 10). Although a purely acoustic amplification cycle has been considered as well (Blondin and Mezzacappa 2006), analytical (Laming 2007, 2008; Yamasaki and Yamada 2007) and numerical (Ohnishi et al. 2006; Scheck et al. 2008; Ferna ´ndez and Thompson 2009a, b) studies have on the whole supported the advective-acoustic cycle, culminating in the work of Guilet and Foglizzo (2012) who sharpened and summarized the arguments in favor of this amplification mechanism.

Different from convection, the SASI is an oscillatory instability with a periodicity T SASI that is set by the sum of the advective and acoustic crossing times s adv and s ac between the shock and the deceleration region (Foglizzo et al. 2007),

Z

Z

<!-- formula-not-decoded -->

The advective time scale usually dominates, and neglecting a weak dependence on the PNS mass, one can determine empirically that the period of the ' ¼ 1 mode of the SASI roughly scales as (Mu ¨ller and Janka 2014),

/C18

/C19

/C18

/C19

<!-- formula-not-decoded -->

where R PNS is the PNS radius. SASI-induced fluctuations in the neutrino emission (Lund et al. 2010; Tamborra et al. 2013; Mu ¨ller and Janka 2014; Mu ¨ller et al. 2019) and gravitational waves (Kuroda et al. 2016a, 2018; Andresen et al. 2017) could provide direct observational confirmation for the SASI if this frequency can be identified in spectrograms of the neutrino or gravitational wave signal.

The growth rate of the SASI is set both by the period T SASI and the quality factor Q of the amplification cycle (Foglizzo et al. 2006, 2007),

<!-- formula-not-decoded -->

where Q depends on the coupling between vortical and acoustic waves at the shock and in the deceleration region, and hence on the details of the density profile and the thermodynamic stratification. Nuclear dissociation and recombination also affect the SASI growth rate and saturation amplitude (Ferna ´ndez and Thompson 2009a, b).

Interplay of SASI, convection, and neutrino heating In reality, the SASI grows in an accretion flow with neutrino heating, and in 2D, it is not trivial at first glance to

<!-- image -->

distinguish SASI and convection in the non-linear phase where both instabilities lead to a similar ' ¼ 1 ''sloshing'' flow. Nonetheless, a clear distinction between a SASI- and convection-dominated regime already emerged in 2D models using gray (Scheck et al. 2008) or multi-group (Mu ¨ller et al. 2012a) neutrino transport, or simpler light-bulb models (Ferna ´ndez et al. 2014): Different from convectiondominated models SASI-dominated models clearly show an oscillatory growth of the multipole coefficients of the shock surface and coherent wave patterns in the post-shock cavity in the linear regime, and maintain a rather clear quasi-periodicity even in the non-linear regime. The distinction between the two different regimes tends to become more blurred around shock revival, when large-scale convective modes emerge and the periodicity of the SASI oscillations is eventually broken.

The criterion v /C25 3 roughly separates the two regimes even though unstable SASI modes can in principle exist above this value. The reason is likely that convection destroys the coherence of the waves involved in the SASI amplification cycle (Guilet et al. 2010) if v [ 3. For v \ 3, high quality factors ln j Q j /C24 2 can be reached and result in rapid SASI growth. In terms of PNS parameters and progenitor parameters, such low values of v \ 3 are encountered in case of rapid PNS and shock contraction (Scheck et al. 2008) and appear to occur preferentially in high-mass progenitors with high mass accretion rates (Mu ¨ller et al. 2012b), although a detailed survey of the progenitor-dependence of the v -parameter is still lacking.

Three-dimensional simulations with neutrino transport (Hanke et al. 2013; Tamborra et al. 2014b; Kuroda et al. 2016a; Mu ¨ller et al. 2017a; Ott et al. 2018; O'Connor and Couch 2018a) as well as simplified leakage and light-bulb models (Couch and O'Connor 2014; Ferna ´ndez 2015) show an even cleaner distinction between the SASI- and convection-dominated regimes for several reasons. The convective eddies remain smaller in the non-linear stage than in 2D because of the forward cascade, and without the constraint of axisymmetry, the convective flow is not prone to artificial oscillatory sloshing motions. The SASI, on the other hand, exhibits a cleaner periodicity prior to shock revival in 3D, and can develop a spiral mode that is very distinct from convective flow (e.g., Blondin and Mezzacappa 2007; Ferna ´ndez 2010; Hanke et al. 2013; see also Sect. 5.3 for possible implications on neutron star birth periods). Self-consistent models show that the post-shock flow can transition back and forth between the convection- and SASIdominated regime as the accretion rate and PNS parameters, and hence the v -parameter change (Hanke et al. 2013).

Saturation mechanism Guilet et al. (2010) argued that parasitic Kelvin-Helmholtz and Rayleigh-Taylor instabilities are responsible for the non-linear saturation of the SASI, and showed that this mechanism can explain the saturation amplitudes in the adiabatic simulations of Ferna ´ndez and Thompson (2009b). Assuming that the Kelvin-Helmholtz instability is the dominant parasitic mode in 3D, one can derive (Mu ¨ller 2016) a scaling law for the turbulent velocity fluctuations d v in the saturated state,

<!-- image -->

<!-- formula-not-decoded -->

which is in good agreement with self-consistent 3D simulations. Interestingly, this scaling results in similar turbulent velocities as in the convection-dominated regime for conditions typically encountered in supernova core (Mu ¨ller 2016).

The saturation of the SASI can also be understood as a self-adjustment to marginal stability (Ferna ´ndez et al. 2014), which is a closely related concept. As the SASI grows in amplitude, the flow is driven towards h v i /C25 3, but stays slightly below this critical value (Ferna ´ndez et al. 2014).

Effect on shock revival The SASI provides similar beneficial effects as convection to increase the shock radius and bring the accretion flow closer to a neutrino-driven runaway, i.e., it generates turbulent pressure, brings high-entropy bubbles to large radii, channels cold matter towards the PNS, and converts turbulent kinetic energy thermal energy throughout the gain region by turbulent dissipation. Due to the different instability mechanism (which feeds on the energy of the accretion flow directly instead of the neutrino energy deposition), and the different flow pattern (which affects the rate of turbulent dissipation), the quantitative effect on shock revival can be different from convection. Using light-bulb simulations Ferna ´ndez (2015) indeed found a significantly lower critical luminosity in the SASI-dominated regime than in the convection-dominated regime and a lower critical luminosity in 3D by /C24 20 % compared to 2D, which he ascribed to the ability of the spiral mode to store more kinetic energy than sloshing modes in 2D. An even bigger difference to convective models (albeit with a different and very idealized setup) was found by Cardall and Budiardja (2015). Self-consistent simulations, on the other hand, have not found higher explodability in 3D in the SASI-dominated regime (Melson et al. 2015a). The reason for this discrepancy could, e.g., lie in the feedback of shock expansion on neutrino heating, but is not fully understood at this stage. Cardall and Budiardja (2015) also observed considerably more stochastic variations in shock revival in the SASI-dominated regime in their idealized models (i.e., a smeared-out critical luminosity threshold), but it again remains to be seen whether this is borne out by self-consistent 3D models, where the SASI oscillations tend to be of smaller amplitude and shorter period than in Cardall and Budiardja (2015).

## 4.5 Perturbation-aided explosions

Progenitor asphericities from convective shell burning can aid shock revival by affecting both the growth and saturation of convection of the SASI. That a higher level of seed perturbations leads to a faster growth of non-radial instabilities behind the shock and thereby fosters explosions (as in the early studies of Couch and Ott 2013; Couch et al. 2015) may be intuitive, but appears less important in practice. In self-consistent simulations, shock revival typically occurs only once convection or the SASI have already reached the stage of non-linear saturation, and it is rather the permanent ' 'forcing'' by infalling perturbations that matters (Mu ¨ller and Janka 2015; Mu ¨ller et al. 2017a). In either case, it is useful to separately consider (a) how

<!-- image -->

the initial perturbations in the porgenitor are translated to perturbations ahead of the shock, and (b) how the infalling perturbations interact with the shock and the postshock flow.

Initial state and infall phase Typically, the Si and O shell (and sometimes a Ne shell) are the only active convective shells that can reach the shock at a sufficiently early post-bounce time to affect shock revival. As described in Sect. 3.2, these shells are characterized by Mach numbers Maprog /C24 0 : 1 with significant variations between different shells and progenitors, and can have a wide range of dominant angular wave numbers ' . Due to its subsonic nature, the flow is almost solenoidal with r/C1 ð q v Þ /C25 0, and density perturbations dq = q /C24 Ma 2 prog are small within convective zones. Viewed as a superposition of linear waves, the convective flow consists mostly of vorticity and entropy waves with little contribution from acoustic waves.

From analytic studies of perturbed Bondi accretion flows in the limit r ! 0 in a broader context (Kovalenko and Eremin 1998; Lai and Goldreich 2000; Foglizzo and Tagger 2000), it is known that such initial perturbations are amplified during infall, and that acoustic waves are generated from the vorticity and entropy perturbations. Estimating the pre-shock perturbations for the problem at hand (Takahashi and Yamada 2014; Mu ¨ller and Janka 2015; Abdikamalov and Foglizzo 2020) involves some subtle differences, but the upshot is rather simple: Advectiveacoustic coupling generates strong acoustic perturbations ahead of the shock that scale linearly with the convective Mach number at the pre-collapse stage (Mu ¨ller and Janka 2015; Abdikamalov and Foglizzo 2020),

<!-- formula-not-decoded -->

According to simulations (Mu ¨ller et al. 2017a) and analytic theory (Abdikamalov and Foglizzo 2020), this scaling is roughly independent of the wave number ' . 12

Shock-turbulence interaction and forced shock deformation The infalling perturbations affect the shock and the post-shock flow in several ways (Mu ¨ller and Janka 2015; Mu ¨ller et al. 2016b). They provide a continuous flux of acoustic and tranverse kinetic energy into the gain region, and also create post-shock density perturbations that will be converted into turbulent kinetic energy by buoyancy. Moreover, the shock becomes deformed due to the anisotropic ram pressure (Fig. 11), which results in fast lateral flow behind the shock, i.e., in the generation of additional transverse kinetic energy. Thus, more violent turbulent flow can be maintained in the gain region, which is conducive to shock revival [cf. Eq. (49)]. If the infalling perturbations are of large scale, the deformation of the shock creates large and stable high-entropy bubbles. This is also helpful for shock revival since runaway shock expansion in multi-D appears to require the formation of such large bubbles

12 If strong acoustic perturbations were present at the pre-collapse stage, these modes with higher ' would grow faster during the linear stage (Takahashi and Yamada 2014), but quickly undergo non-linear damping (Mu ¨ller and Janka 2015).

<!-- image -->

Fig. 11 Interaction of infalling perturbations with the shock and the post-shock flow, illustrated by snapshots of the entropy (in units of k b = nucleon, left panel) and the absolute value of the non-radial velocity (in units of km s /C0 1 , right panel) in the 12 : 5 M /C12 model of Mu ¨ller et al. (2019) at a post-bounce time of 510ms. The left panel also shows the deformation of the isodensity surface with q ¼ 7 /C2 10 6 g cm /C0 6 (red curve). Due to the infalling density perturbations, the pre-shock ram pressure is anisotropic and creates a protrusion of the shock. Additional energy is pumped into non-radial motions in the gain region both because of substantial lateral velocity perturbations ahead of the shock and because of the oblique infall of material through the deformed shock

<!-- image -->

with sufficient buoyancy to rise and expand against the supersonic drag of the infalling material (Ferna ´ndez et al. 2014; Ferna ´ndez 2015).

There is as yet no comprehensive quantitative theory for the interaction of infalling perturbations with the shock and the instabilities in the gain region, but several studies have investigated aspects of the problem. Using order-of-magntitude estimates, Mu ¨ller et al. (2016b) argued that turbulent motions are primarily boosted by the action of buoyancy on the injected post-shock density perturbations. This hypothesis is supported by controlled parameterized simulations of shockturbulence interactions in planar geometry (Kazeroni and Abdikamalov 2019). Mu ¨ller et al. (2016b) also attempted to derive a correction for the saturation value of turbulent kinetic energy depending on the convective Mach number Maprog and wave number ' in the progenitor. They predicted a reduction of the critical luminosity functional L crit ¼ ð L m E 2 m Þ crit by

<!-- formula-not-decoded -->

in terms of the heating efficiency g heat and the accretion efficiency g heat ¼ L = ð GM \_ M = R g Þ . However, the analysis of Mu ¨ller et al. (2016b) did not account in detail for the interaction of the infalling perturbations with the shock. This has been investigated using linear perturbation theory (Takahashi et al. 2016; Abdikamalov et al. 2016, 2018; Huete et al. 2018; Huete and Abdikamalov 2019). As a downside, this perturbative approach cannot easily capture the non-linear interaction of the injected perturbations with fully developed neutrino-driven convection and the SASI, but Huete et al. (2018) recently incorporated the effects of

<!-- image -->

buoyancy downstream of the shock. The more sophisticated treatment of Huete et al. (2018) predicts a similar effect size as Eq. (63).

Phenomenology of perturbation-aided explosions Whatever its theoretical justification, Eq. (63) successfully captures trends seen in 2D and 3D simulations of perturbation-aided explosions starting from parameterized initial conditions or 3D progenitor models. Both high Mach numbers J 0 : 1 and large-scale convection with ' . 4 are required for a significant beneficial effect on the heating conditions (Mu ¨ller and Janka 2015). In this case, the perturbations can be the decisive factor for shock revival as in the 18 M /C12 model of Mu ¨ller et al. (2017a). In leakage-based models with high heating efficiency early on (Couch and Ott 2013; Couch et al. 2015), the effect is smaller, especially if the pre-collapse asphericities are restricted to medium-scale modes as in octant simulations (Couch et al. 2015).

By now, there is a handful of exploding supernova models that use multi-group neutrino transport and 3D progenitor models (Mu ¨ller et al. 2017a, 2019). While this is encouraging, more 3D simulations are needed to determine to what extent convective seed perturbations generally contribute to robust explosions. At present, one can nonetheless extrapolate the effect size based on the properties of convective shells in 1D stellar evolution models using Eq. (63). Analysing over 2000 supernova progenitors computed with the KEPLER code Collins et al. (2018) predict a substantial reduction of the critical luminosity due to perturbation by 10% or more in the mass range between 15 M /C12 and 27 M /C12 , and in isolated low-mass progenitors. Below 15 M /C12 , the expected reduction is usually 5% or less, which could still make the convective perturbations one of several important ingredients for robust explosions. In the vast majority of progenitors, only asphericities from oxygen shell burning are expected to have an important dynamic effect.

## 4.6 Outlook: rotation and magnetic fields in neutrino-driven explosions

Earlier on, we already briefly touched simulations of magnetorotational explosion scenarios and the uncertainties that still beset this mechanism. It is noteworthy that rotation and magnetic fields could also play a role within the neutrino-driven paradigm.

Rotationally-supported explosions Since early attempts to study the impact of rotation on neutrino-driven explosions either employed a simplified neutrino treatment (e.g., Kotake et al. 2003; Fryer and Warren 2004; Nakamura et al. 2014; Iwakami et al. 2014) or were restricted to 2D in the case of models with multi-group transport (Walder et al. 2005; Marek and Janka 2009; Suwa et al. 2010), more robust conclusions had to wait for 3D simulations with multi-group neutrino transport (Takiwaki et al. 2016; Janka et al. 2016; Summa et al. 2018). The 3D simulations indicate that the overall effect of rapid rotation is to support neutrinodriven explosions. Centrifugal support reduces the infall velocities and hence the average ram pressure at the shock (Walder et al. 2005; Janka et al. 2016; Summa et al. 2018). Moreover, 3D neutrino hydrodynamics simulations of rotating models tend to develop a strong spiral SASI (Janka et al. 2016; Summa et al. 2018). This is

<!-- image -->

in line with analytic theory (Yamasaki and Foglizzo 2008) and idealized simulations (Iwakami et al. 2009; Blondin et al. 2017; Kazeroni et al. 2018), which demonstrated that rotation enhances the growth rate of the prograde spiral mode and stabilises the retrograde mode. For sufficiently rapid rotation, an even more violent spiral corotation instability can occur (Takiwaki et al. 2016). There is also a subdominant adverse effect, since lower neutrino luminosities and mean energies at low latitudes close to the equatorial plane are detrimental for shock revival (Walder et al. 2005; Marek and Janka 2009; Summa et al. 2018), which is particularly relevant since the explosion tends to be aligned with the equatorial plane in the case of rapid rotation (Nakamura et al. 2014). Summa et al. (2018) found that the overall combination of these effects can be encapsulated by a further modification of the critical luminosity,

s

ffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffi ffi

<!-- formula-not-decoded -->

Here j is the spherically-averaged angular momentum of the shell currently falling through the shock. The last factor accounts for the reduced pre-shock velocities, and the effect of stronger non-radial flow in the gain region is implicitly (but not predictively) accounted for in the turbulent Mach number.

Magnetic fields without rotation Given the expected pre-collapse spin rates, rotation is unlikely to have a major impact in the vast majority of supernova explosions. It is harder to exclude a significant role of magnetic fields a priori . Even if the progenitor core rotates slowly and does not have strong magnetic fields, convection and the SASI might furnish some kind of turbulent dynamo process that could generate dynamically relevant fields in the gain region. There could also be other processes to provide dynamically relevant magnetic fields, e.g., the accumulation of Alfve ´n waves at an Alfve ´n surface (Guilet et al. 2011) or the injection of Alfve ´n waves generated in the PNS convection zone (Suzuki et al. 2008).

The simulations available so far do not suggest that sufficiently high field strengths can be reached by a small-scale turbulent dynamo. In idealized 2D and 3D simulations of Endeve et al. (2010, 2012), the SASI indeed drives a small-scale turbulent dynamo, and strong field amplification occurs locally up to equipartition and super-equipartition field strengths, especially when a strong spiral mode develops. On larger scales, the magnetic field energy remains well below equipartition, however, and does not become dynamically important. The total magnetic energy in the gain region remains one order of magnitude smaller than the turbulent kinetic energy, and the field does not organize itself into large-scale structures. The situation is similar in the 2D neutrino hydrodynamics simulations of non-rotating progenitors of Obergaulinger et al. (2014) for initial field strengths of up to 10 11 G, with even lower ratios between the total magnetic and turbulent kinetic energy in the gain region. Only for initial field strengths of /C24 10 12 G, which yields magnetar-strength fields after collapse, do Obergaulinger et al. (2014) find that magnetic fields become dynamically important and accelerate shock revival. If the

<!-- image -->

fossil field hypothesis for magnetars is correct and the fields of the most strongly magnetized main-sequence stars translate directly to supernova progenitor and neutron star fields by flux conservation (Ferrario and Wickramasinghe 2006), such conditions for magnetically-aided explosions might be still realized in nature in a substantial fraction of core-collapse events.

Ultimately, a thorough exploration of resolution effects and initial field configurations in the convection- and SASI-dominated regime will be required in 3D to confidently exclude a major role of magnetic field in weakly magnetized, slowly rotating progenitors. First tentative results from 3D MHD simulations with neutrino transport (Fig. 12) suggest a picture of fibril flux concentrations with equipartition field strengths, and sub-equipartition fields in most of the volume, akin to the situation in solar convection (e.g., Solanki et al. 2006).

## 4.7 Proto-neutron star convection and LESA instability

Prompt convection Convective instability also develops inside the PNS. As already recognized in the late 1980s (Bethe et al. 1987; Bethe 1990), a first episode of ' 'prompt convection'' occurs within milliseconds after bounce around a mass coordinate of /C24 0 : 8 M /C12 as the shock weakens and a negative entropy gradient is established. The negative entropy gradient is, however, quickly smoothed out, and the convective overturn has no bearing on the explosion mechanism, although it can leave a prominent signal in gravitational waves (see reviews on the subject; Ott 2009; Kotake 2013; Kalogera et al. 2019).

<!-- image -->

Fig. 12 Entropy s in k b = nucleon (left panel) and the logarithm log 10 P B = P gas of the ratio between the magnetic pressure P B and the gas pressure P gas (right panel) in a 3D simulation (left half of panels) and a 2D simulation (right half of panels) of the slowly rotating progenitor 15 M /C12 progenitor m15b6 of Heger et al. (2005) with the COCONUT-FMT code. The initial field is assumed to be combination of a dipolar poloidal field and a toroidal field. Outside convective zones, the field strength is taken from the progenitor, inside convective zones, the magnetic pressure is set to a fraction of 10 /C0 4 of the thermal pressure. The figures shows meridional slices 140ms after bounce. Field amplification is driven by convection. Strong fields are generated in regions of strong shear, but these strong field are highly localized, and the total magnetic energy in the gain region remains much smaller than the turbulent kinetic energy and thermal energy

<!-- image -->

Proto-neutron star convection Convection inside the PNS is triggered again latter as neutrino cooling establishes unstable lepton number (Epstein 1979) and entropy gradients (see profiles in Fig. 9). PNS convection was investigated extensively in the 1980s and 1990s as a possible means of enhancing the neutrino emission from the PNS, which would boost the neutrino heating and thereby aid shock revival (e.g., Burrows 1987; Burrows and Lattimer 1988; Wilson and Mayle 1988, 1993; Janka and Mu ¨ller 1995; Keil et al. 1996). In particular, Wilson and Mayle (1988, 1993) assumed that PNS convection operates as double-diffusive ''neutron finger'' instability that significantly increases the neutrino luminosity.

None of the modern studies of PNS convection since the mid-1990s (Keil et al. 1996; Buras et al. 2006a; Dessart et al. 2006) found a sufficiently strong effect of PNS convection on the neutrino emission for a significant impact on shock revival. PNS convection indeed increases the heavy flavor neutrino luminosity by /C24 20 % at post-bounce times of J 150 ms, leaves the electron neutrino luminosity about equal, but tends to decrease the electron antineutrino luminosity, and reduces the mean energy of all neutrino flavors (Buras et al. 2006a). This can be explained by the effects of PNS convection on the bulk structure of the PNS, namely a modest increase of the PNS radius and a higher electron fraction (due to mixing) close to the neutrinosphere of m e and /C22 m e (Buras et al. 2006a). Convective instability appears to be governed by the usual Ledoux criterion and does not develop as a double-diffusive instability in the simulations.

LESA instability Although PNS convection does not have a decisive influence on shock revival, its indirect effect on the gain region is quantitatively important; effectively PNS convection changes the inner boundary condition for the flow in the gain region. PNS convection also has an important impact on the neutrino signal from the PNS cooling phase (e.g., Roberts et al. 2012; Mirizzi et al. 2016), and may provide a sizable contribution to the gravitational wave signal (Marek et al. 2009; Yakunin et al. 2010; Mu ¨ller et al. 2013; Andresen et al. 2017; Morozova et al. 2018).

Moreover, starting with (Tamborra et al. 2014a), the dynamics of PNS convection has proved more intricate upon closer inspection in recent years with potential repercussions on nucleosynthesis and gravitational wave emission. In their 3D simulations, Tamborra et al. (2014a) noted that a pronounced ' ¼ 1 asymmetry in the electron fraction develops in the PNS convection zone, which leads to a sizable anisotropy in the radiated lepton number flux, i.e., the fluxes of electron neutrinos and antineutrinos show a dipole asymmetry with opposite directions. The direction of the dipole varies remarkably slowly compared to the characteristic time scales of the PNS. Curiously, no such pronounced dipole was seen in the velocity field, which remained dominated by small-scale eddies. This is very different from convection in the gain region or the convective burning in the progenitors, where the asymmetries in the entropy and the composition are reflected in the velocity field as well. This unusual phenomenon, which is illustrated in Fig. 13, has been christened LESA (''Lepton Number Emission Self-Sustained Asymmetry''), and could have important repercussions on the composition of the ejected matter, whose Y e is sensitive to the differences in electron neutrino and antineutrino emission.

<!-- image -->

Fig. 13 LESA instability in a simulation of an 18 M /C12 star at time of 453 ms after bounce, illustrated by 2D slices showing the electron fraction Y e (left) and the radial velocity in units of km s /C0 1 in the PNS convection zone. Note that the Y e distribution in the PNS convection zone between radii of 10 km and 20 km shows a clear dipolar asymmetry, whereas the radial velocity field is dominated by small-scale modes superimposed over a much weaker dipole mode. Image repoduced with permission from Powell and Mu ¨ller (2019), copyright by the authors

<!-- image -->

Since then this phenomenon has been reproduced by many 3D simulations using very different methods for neutrino transport (O'Connor and Couch 2018a; Janka et al. 2016; Glas et al. 2019; Powell and Mu ¨ller 2019; Vartanyan et al. 2019a), and even in the 3D leakage models of Couch and O'Connor (2014) The dipolar Y e asymmetry has also been seen in the 2D Boltzmann simulations of Nagakura et al. (2019). This demonstrates that the LESA is a robust phenomenon; claims that it depends on the details of neutrino transport (Nagakura et al. 2019) are not convincing.

The nature of this instabiity is still not fully understood. Tamborra et al. (2014a) initially suggested a feedback cycle between ' ¼ 1 shock deformation, a dipolar asymmetry in the accretion flow, and a dipolar asymmetry in the lepton fraction in the PNS convection zone. However, more recent studies suggest that the LESA does not depend on an external feedback cycle between asymmetries in the accretion flow and asymmetries in the PNS convection zone. Glas et al. (2019) demonstrated that LESA can be even more pronounced in exploding low-mass progenitor models with low accretion rates, which suggests that the mechanism behind LESA works within the PNS convection zone.

This, however, leaves the question why the flow within the PNS convection zone would organise itself to generate a dipolar lepton fraction asymmetry. Some papers have, however, formulated first qualitative arguments to suggest that there is an internal mechanism for a dipole asymmetry in the lepton fraction, and that LESA may just be a very peculiar manifestation of buoyancy-driven convection. What appears to play a role is that the lepton fraction gradient becomes stabilizing against convection in the middle of the PNS convection zone (Janka et al. 2016; Powell and Mu ¨ller 2019). Janka et al. (2016) suggested that this can give rise to a positive feedback loop because a hemispheric lepton asymmetry will attenuate or enhance the stabilizing effect in the different hemispheres, thereby leading to more vigorous

<!-- image -->

convection in one hemisphere, which in turn maintains the lepton asymmetry. On a different note, Glas et al. (2019) sought to explain the large-scale nature of the asymmetry by applying the concept of a critical Rayleigh number for thermallydriven convection (Chandrasekhar 1961). However, one still needs to account for the fact that the typical scales of the velocity and lepton number perturbations appear remarkably different in the PNS convection zone. This was confirmed by the quantitative analysis of Powell and Mu ¨ller (2019), who found a very broad turbulent velocity spectrum peaking around ' ¼ 20, which conforms neither to Kolmogorov or Bolgiano-Obukhov scaling for stratified turbulence. Powell and Mu ¨ller (2019) suggested that this could be explained by the scale-dependent effective buoyancy experienced by eddies of different sizes as they move across the partially stabilized central region of the PNS convection zone. Powell and Mu ¨ller (2019) also remarked that the non-linear state of PNS convection is characterized by a balance between the convective and diffusive lepton number flux. All these aspects suggest that the LESA could be no more than a manifestation of PNS convection, but that PNS convection is in fact quite dissimilar from the high-Pe ´clet number convection as familiar from the gain region or the late convective burning stages. A satisfactory explanation of the phenomenon likely needs to go beyond concepts from linear stability theory and the usual global balance arguments behind MLT, and will have to take into account scale-dependent forcing and dissipation, and the ' 'doubleradiative'' nature of the instability.

Since the stabilising lepton number gradient in the middle of the PNS convection zone figures prominently in these attempts to understand LESA, one might justifiably ask whether there is some role for double-diffusive instabilities in the PNS after all. Local stability analysis (Bruenn et al. 1995; Bruenn and Dineva 1996; Bruenn et al. 2004) in fact suggests that double-diffusive instabilities (termed leptoentropy fingers and lepto-entropy semiconvection) should occur in the PNS. But why were such double-diffusive instabilities never identified in multi-D simulations so far? Further careful analysis and interpretation of the simulation results and theory is in order to clarify this. One possible interpretation could be that the characteristic step-like lepton number profile established by LESA (Powell and Mu ¨ller 2019) is actually a manifestation of layer formation in the subcritical regime as familiar from semiconvection (Proctor 1981; Spruit 2013; Garaud 2018). However, the slow, global turnover motions in LESA do not readily fit into this picture. 13 One should also beware premature conclusions because PNS convection is an inherently difficult regime for numerical simulations due to small convective Mach numbers of order /C24 0 : 01 and the importance of diffusive effects. The potential issues go beyond the question of resolution and unphysically high numrical Reynolds numbers (cf. 2.1.2), and there are concrete reasons to investigate these in more depth. For example, although different codes agree qualitatively

13 That the LESA is also seen in models without lateral diffusion may not an obstacle for this interpretation. Lateral diffusion is essential to obtain semiconvective overstability, but layer formation can occur below the threshold for overstability (Proctor 1981; Spruit 2013).

<!-- image -->

concerning the region of instability and the qualitative features of the convective flow, substantial differences in the turbulent kinetic energy density have been reported in a comparison between the ALCAR and VERTEX codes in the PNS convection zone, even though the agreement between the codes is otherwise excellent (Just et al. 2018). While it is unlikely that the uncertainties in models PNS convection have any impact on the problem of shock revival, they need to be addressed to obtain a better understanding of LESA and reliable predictions of gravitational wave signals and the nucleosynthesis conditions in the neutrino-heated ejecta.

## 5 The explosion phase

Regardless of whether the explosion is driven by neutrinos or magnetic fields, there is no abrupt transition to a quasi-spherical outflow after shock revival. In this section, we shall focus on the situation in neutrino-driven explosions, which has already been quite thoroughly explored.

## 5.1 The early explosion phase

In typical neutrino-driven models, the multi-dimensional flow structure in the early explosion phase appears qualitatively similar to the pre-explosion phase at first glance. Buoyancy-driven outflows and accretion downflows persist for hundreds of milliseconds to seconds and allow for simultaneous mass accretion and ejection. Because of the ongoing accretion, high neutrino luminosities and hence high heating rates can be maintained to continually dump energy into the developing explosion. As the shock radius slowly increases, large-scale ' ¼ 1 and ' ¼ 2 modes start to dominate the flow irrespective of whether medium-scale convection or large-scale SASI modes dominated prior to shock revival, even though 2D explosion models probably tended to exaggerate this effect. The basic features of this pictures have held since the 1990s (Herant et al. 1992; Shimizu et al. 1993; Yamada et al. 1993; Janka et al. 1993; Herant et al. 1994; Burrows et al. 1995; Janka and Mu ¨ller 1995, 1996), and have proved critical for explaining the energetics of core-collapse supernovae. 14 Even in electron-capture supernova progenitors, which explode even without the help of multi-dimensional effects (Kitaura et al. 2006), there is a brief phase of convective overturn after shock revival (Wanajo et al. 2011).

More recent 3D explosion models using multi-group transport (Takiwaki et al. 2014; Melson et al. 2015a; Lentz et al. 2015; Mu ¨ller 2015; Mu ¨ller et al. 2017a, 2019; Burrows et al. 2020) have confirmed this picture, but paved the way towards a more quantitative theory of the explosion phase. In massive progenitors, shock expansion is usually sufficiently slow for one or two dominant bubbles of neutrino-heated ejecta to form (Fig. 14). Only at the low-mass end of the

14 Critiques of the neutrino-driven mechanism have occasionally overlooked (Papish et al. 2015) and then ultimately rebranded the simultaneous outflows and downflows as ''jittering jets'' (Soker 2019). In this latest instalment, the alternative jittering-jet scenario seems to have come down to little more than a question of unconventional terminology for well-established phenomena in neutrino-driven explosions.

<!-- image -->

Fig. 14 Volume renderings of the entropy in different 3D supernova simulations showing the emergence of stable large-scale plumes around and after shock revival as a common phenomenon despite differences in resolution and in the neutrino transport treatment. The outer translucent surface is the shock, the structures inside are neutrino-heated high-entropy bubbles: a 15 M /C12 model of Lentz et al. (2015) with a unipolar explosion geometry at a post-bounce time of 400 ms. b 3 M /C12 He star model of Mu ¨ller et al. (2019) at 1238 ms, with two prominent plumes in the 7 o'clock and 11 o'clock directions and weaker shock expansion on the opposite side. c 20 M /C12 model of Burrows et al. (2020, Fig. 8) with a more dipolar explosion geometry at 651 ms. d 11 : 2 M /C12 model of Nakamura et al. (2019, Fig. 6) at a time of 991 ms. Images reproduced with permission from [a] Lentz et al. (2015), copyright by AAS; [c] from Burrows et al. (2020) and [d] from Nakamura et al. (2019), copyright by the authors

<!-- image -->

progenitor spectrum (Melson et al. 2015b; Gessner and Janka 2018) do the convective structures freeze out so quickly that the neutrino-heated ejecta are organized in medium-scale bubbles instead of a unipolar or bipolar structure.

The detailed dynamics of the outflows and downflows proved to be significantly different in 3D compared to 2D, and that only restricted insights on explosion and remnant properties and nucleosynthesis can be gained from the impressive corpus of successful 2D simulations with multi-group transport (Buras et al. 2006a; Marek and Janka 2009; Mu ¨ller et al. 2012a, b, 2013; Janka 2012; Janka et al. 2012; Suwa et al. 2010, 2013; Bruenn et al. 2013, 2016; Nakamura et al. 2015; Burrows et al. 2018; Pan et al. 2018; O'Connor and Couch 2018b). Except at the lowest masses, the 2D simulations are uniformly characterized by almost unabated accretion through fast downflows that reach directly to the bottom of the gain region, by

<!-- image -->

outflows that are often weak and intermittent, and by a halting rise of explosion energies. Long-time 2D simulations showed that this situation can persist out to more than 10 s (Mu ¨ller 2015), and as a result, implausibly high neutron star masses are reached. The halting growth of explosion energies in 2D can partly be explained by the topology of the flow which lends itself to outflow constriction by equatorial downflows, but the primary difference between 2D and 3D lies in the velocity of the downflows. Melson et al. (2015b) already noticed that the downflows appear to subside more quickly in their 3D model of a low-mass progenitor, which they ascribed to the forward turbulent cascade in 3D; this led to a slight enhancement of the explosion energy by 10% in 3D compared to 2D. In more massive progenitors with stronger accretion after shock revival stronger braking of the downflows in 3D compared to 2D is even more evident (Mu ¨ller 2015). Instead of crashing into a secondary accretion shock at /C24 100 km at a sizable fraction of the free-fall velocity, the downflows are gently decelerated, and secondary shocks rarely form. Mu ¨ller (2015) ascribed this pathology of the 2D models to the behavior of the KelvinHelmholtz instability between the outflows and downflows, which is stabilized at high Mach numbers in 2D, but can always grow in 3D (Gerwin 1968). Since the typical Mach number of the downflows is higher during the explosion phase, the assumption of 2D symmetry becomes even more problematic than during the accretion phase.

## 5.2 Explosion energetics

Estimators for the explosion energy Strictly speaking, the final demarcation between ejected and accreted material 15 cannot be determined before the explosion becomes kinetically dominated after shock breakout, and the same holds true for the final explosion energy E exp . It is, however, customary and useful to consider the diagnostic explosion energy E diag (often shortened to ' 'diagnostic energy'' or ' 'explosion energy'' when there is no ambiguity), which is defined as the total energy of the material that is nominally unbound at any given instance (Buras et al. 2006a; Mu ¨ller et al. 2012b; Bruenn et al. 2013). By definition the diagnostic energy will eventually asymptote to E exp, but might do so only over considerably longer time scales than can be simulated with neutrino transport. In particular, the diagnostic energy can in principle decrease as the shock sweeps up bound material from the outer shells. To account for this, one can correct E diag for the binding energy of the material ahead of the shock (' 'overburden'') to obtain a more conservative estimate for E exp (Bruenn et al. 2016). In practice, E diag usually rises monotonically because energy continues to be pumped into the ejecta over seconds, but there are exceptions, most notably in cases of early black hole formation (Chan

15 We avoid the term ''mass cut'', which is commonly used for describing artificial 1D explosion models. The boundary of the ejecta region is not a sphere, and does not correspond to a unique mass shell under realistic conditions.

<!-- image -->

et al. 2018). In most cases, one expects that E diag levels out after a few seconds and then provides a good estimate for E exp .

Explosion energies from self-consistent simulations Unfortunately, E diag has not levelled off in most of the available self-consistent 3D explosion models, though the growth of the explosion energy has already slowed down significantly in some longtime simulations using the COCONUT-FMT code (Mu ¨ller et al. 2017a, 2018, 2019). Even in 2D, only some of the models of the Oakridge group appear to have approached their final explosion energy (Bruenn et al. 2016).

This means that no final verdict on the fidelity of the simulations can be pronounced based on a comparison with observationally inferred explosion energies. The models of Mu ¨ller et al. (2017a, 2018, 2019) and Bruenn et al. (2016), whose explosion energies are admittedly on the high side among modern simulations, have demonstrated that neutrino-driven explosions can reach energies of up to 8 /C2 10 50 erg. Similarly, plausible nickel masses of several 0 : 01 M /C12 appear within reach, although no firm statements can be made for COCONUT-FMT models due to uncertainties in the Y e of the ejecta from the approximative transport treatment, and due to the highly simplified treatment of nucleon recombination. Explosion energies beyond 10 51 erg may simply be a matter of longer simulations, different progenitor models, and slightly improved physics; and there may be no conflict with the distribution of observationally inferred explosion energies (Kasen and Woosley 2009; Pejcha and Prieto 2015; Mu ¨ller et al. 2017b) of Type IIP supernovae. First attempts to extrapolate the non-converged explosion energies from simulations and compare them to observations using a rigorous statistical framework (Murphy et al. 2019) indicate that the predicted values are still somewhat too low, but Murphy et al. (2019) also point out that conclusions are premature due to biases and uncertainties in the comparison.

Growth of the explosion energy Even at this stage, the simulations already elucidate how the energy of neutrino-driven explosions is determined. Upon closer inspection, the energy budget of the ejecta is quite complicated and includes contributions from the injection of neutrino-heated material from below, form nucleon recombination and nuclear burning, from the accumulation of bound material by the shock, and from turbulent mixing with the downflows (for a broader discussion, see Marek and Janka 2009; Mu ¨ller 2015; Bruenn et al. 2016). Nonetheless, a few key findings have emerged. The most critical determinant for the growth of E diag is the mass outflow rate \_ M out of neutrino-heated material. Neutrino heating only marginally unbinds the material, and the net contribution to E diag comes from the energy /C15 rec released by nucleon recombination, which occurs at a radius of about 300 km. To first order, the resulting growth rate of the diagnostics energy is (Scheck et al. 2006; Melson et al. 2015b; Mu ¨ller 2015),

<!-- formula-not-decoded -->

In principle, 8 : 8 MeV per nucleon can be released from full recombination to the iron group, but for the relevant entropies and expansion time scales, recombination

<!-- image -->

is incomplete and does not convert all the neutrino-heated ejecta to iron-group elements. Mixing between the outflows and downflows reduces the effective value of e rec further to about 5-6 MeV = nucleon.

The mass outflow rate is roughly determined by the volume-integrated neutrino heating rate and the energy required to lift the material out of the gravitational potential. One can argue that the relevant energy scale is the binding energy at the gain radius j e gain j , so that

<!-- formula-not-decoded -->

with some efficiency parameter g out that accounts for the fact that only part of the neutrino-heated matter is ejected. Initially, one finds g out \ 1 as expected, but Mu ¨ller et al. (2017a) pointed out that the situation changes later after shock revival because much of the ejected material never makes it down to the gain radius and starts is expansion significantly further out with lower initial binding energy. This leads to efficiency parameters g out [ 1, and helps compensate for the declining heating rates as the supply of fresh material to the gain radius slowly subsides.

The situation in 2D models is somewhat different (Mu ¨ller 2015). Here, the ejected material comes from close to the gain radius, and the outflow efficiency g out is lower than in 3D. Although the lack of turbulent mixing results in a higher asymptotic specific total energy and entropy in 2D, the net effect is a slower growth of the explosion energy than in 3D. Moreover, the higher entropies in 2D will result in reduced recombination to the iron group and hence lower nickel masses. Despite these shortcomings, 2D simulations remain of some use because they already allow extensive parameters studies of explodability and explosion and remnant properties (Nakamura et al. 2015).

## 5.3 Compact remnant properties

Accretion rates and remnant masses The forward cascade and the stronger KelvinHelmholtz instabilities between the outflows and downflows in 3D imply that the accretion rate onto the PNS drops more quickly than in 2D (Mu ¨ller 2015). As a result, some self-consistent 3D simulations have been able to determine firm numbers for final neutron star masses (Melson et al. 2015b; Mu ¨ller et al. 2019; Burrows et al. 2019, 2020), barring the possibility of late-time fallback. The predicted neutron star masses appear roughly compatible with the range of observed values (O ¨ zel and Freire 2016; Antoniadis et al. 2016), but as with all other explosion and remnant properties, a robust statistical comparison is not yet possible.

Neutron star kicks Observations show that most neutron stars receive a considerable ' 'kick' ' velocity at birth. The kick velocity is typically a few hundred km s /C0 1 , but there is a broad distribution ranging from very low kicks up to more than 1000 km s /C0 1 (e.g., Hobbs et al. 2005; Faucher-Gigue `re and Kaspi 2006; Ng and Romani 2007). The large-scale ejecta asymmetries that emerge during the explosion provide a possible explanation for this phenomenon (for an overview including

<!-- image -->

other mechanisms such as aniostropic neutrino emission, see Lai et al. 2001; Janka 2017).

The 2D simulations of the 1990s could not yet naturally obtain the full range of observed kick velocities by a hydrodynamic mechanism (Janka and Mu ¨ller 1994), unless unrealistically large seed asymmetries in the progenitor were invoked (Burrows and Hayes 1996). A plausible range of kicks was first obtained in parameterized 2D simulations by Scheck et al. (2004, 2006), thanks to more slowly developing explosions that allowed the ' ¼ 1-mode of the SASI, or an ' ¼ 1 convective mode to emerge. The work of Scheck et al. (2004, 2006) revealed that the kick velocity can grow for well over a second in models with slower shock expansion. They concluded that the kick arises primarily from the asymmetric gravitational pull of overand underdense regions in the ejecta (later termed ' 'gravitational tugboat mechanism'' 16 ) rather than from pressure force and hydrodyanmic momentum fluxes onto the PNS; anisotropic neutrino emission was found to play only a minor role. Subsequent simulations have not fundamentally changed this analysis. Although various studies showed that the momentum flux onto the PNS can be comparable to the gravitational force onto the PNS (Nordhaus et al. 2010a, 2012; Mu ¨ller et al. 2017a), this does not invalidate the tugboat mechanism. Effectively, the contribution of each parcel of accreted material to the PNS momentum via the hydrodynamic flux and the gravitational tug almost cancel, and the net acceleration of the PNS is due to the gravitational pull of the material that is actually ejected.

Three-dimensional simulations have not altered this picture substantially. Even though 2D simulations tend to obtain higher kicks, values of several hundred km s /C0 1 were already obtained in the parameterized 3D simulations of Wongwathanarat et al. (2010b, 2013). Recently, Mu ¨ller et al. (2017a, 2019) performed sufficiently long 3D simulations with multi-group neutrino transport to extrapolate the final kick velocities, which fall nicely within the observed range of up to 1000 km s /C0 1 .

Based on the physics of the kick mechanism, various authors have posited a correlation between the kick and the ejecta mass (Bray and Eldridge 2016) or, using more refined arguments, on the explosion energy (Janka 2017; Vigna-Go ´mez et al. 2018). Tentative support for a loose correlation comes from the small kicks obtained in simulations of low energy, ultra-stripped supernovae (Suwa et al. 2015; Mu ¨ller et al. 2017a) and electron-capture supernovae (Gessner and Janka 2018), and from more recent 3D simulations over a larger range of progenitor masses (Mu ¨ller et al. 2019).

Neutron star spins If the downflows hit the PNS surface with a finite impact parameter, they also impart angular momentum onto the PNS. While this was realized already by Spruit and Phinney (1998), 3D simulations are needed for quantitative predictions of PNS spin-up by asymmetric accretion. The predicted spin-up in 3D models of non-rotating varies. Parameterized simulations (Wongwathanarat et al. 2010b; Rantsiou et al. 2011; Wongwathanarat et al. 2013) tend to find longer neutron star spin periods of hundreds of milliseconds to seconds (but

16 The term ''tugboat mechanism'' was in fact suggested later by Jeremiah Murphy and introduced into the literature in Wongwathanarat et al. (2013).

<!-- image -->

extending down to 100 ms in Wongwathanarat et al. 2013). Recent 3D simulations using multi-group transport (Mu ¨ller et al. 2017a, 2019) obtain spin periods between 20 ms and 2 : 7 s, which roughly coincides with the range of observationally inferred birth periods (Faucher-Gigue `re and Kaspi 2006; Perna et al. 2008; Popov and Turolla 2012; Noutsos et al. 2013). Assuming the core angular momentum is conserved after the collapse to a PNS and not changed by angular momentum transport during the explosion, current stellar evolution models computed with the Tayler-Spruit dynamo, predict spin periods in the same range (Heger et al. 2005), which makes it difficult to draw inferences on the explosion mechanism or the progenitor rotation from the observed spin periods. None of the simulations can as yet explain the spin-kick alignment that is suggested by observations (Johnston et al. 2005; Ng and Romani 2007; Noutsos et al. 2013). Proposed mechanisms for natural spin-kick alignment by purely hydrodyanmic processes (Spruit and Phinney 1998; Janka 2017) have not been borne out by the models. However, the possibility of natural spin-kick alignment in rotating progenitors has yet to be investigated.

Role of the spiral SASI mode In the first idealized simulations of the spiral mode of the SASI, Blondin and Mezzacappa (2007) noted a significant flux of angular momentum into the ''neutron star'' (modeled by an inner boundary condition) that would lead to rapid neutron star rotation in the opposite direction to the SASI flow with angular frequencies of the order of 100 rad s /C0 1 even in the case of non-rotating progenitors. This idea has been explored further in several numerical (Hanke et al. 2013; Kazeroni et al. 2016, 2017) and analytical (Guilet and Ferna ´ndez 2014) studies. The potential for spin-up of non-rotating progenitors may be more modest than initially thought; both numerical and analytical results suggest that the angular momentum imparted onto the PNS is only a few 10 46 erg s, corresponding to spin periods of hundreds of milliseconds (Hanke et al. 2013; Guilet and Ferna ´ndez 2014). Moreover, part of the angular momentum contained in the spiral mode may be accreted after shock revival, negating the separation of angular momentum previously achieved by the SASI. Spin-up and spin-down by SASI in the case of rotating progenitors still merits further investigation; the idealized simulations of Kazeroni et al. (2017) suggest different regimes of random spin-up and spin-down for slow progenitor rotation, systematic spin-down for intermediate rotation, and weaker spin-down for high rotation rates in the regime of the corotation instability. The possibility of magnetic field amplification due to the induced shear in the PNS surface region in the case of significant spin-up or spin-down by the SASI also needs to be explored.

## 5.4 Mixing instabilities in the envelope

Structure of the flow in the later explosion phase As the propagating shock scoops up matter and as the explosion energy levels off, the structure of the post-shock region changes (Fig. 15a). Early on, the post-shock expansion velocities are subsonic and the outflows are accelerated by a positive pressure gradient, but eventually the post-shock flow enters a Sedov-like regime where a positive pressure

<!-- image -->

Fig. 15 a Emergence of Rayleigh-Taylor instability during the propagation of the shock through the envelope, illustrated by spherically averaged profiles of density q , pressure P , and radial velocity vr from a 2D long-time simulation of a 9 : 6 M /C12 star based on the explosion model of Mu ¨ller et al. (2013). At this stage (140 s after the onset of the explosion), the shock has reached the He shell, and a reverse shock has formed. Behind the forward shock, the pressure gradient is positive and decelerates the expansion of the ejecta. Rayleigh-Taylor (RT) instability grows in a region with d q = d r \ 0 behind the shock. Note that the structure of the blast wave can be more complicated in general with several unstable regions and reverse shocks that interact with each other. b Mass fraction isocontours in a 3D model of mixing in SN 1987A. Note that while the biggest Ni-rich Rayleigh-Taylor clumps are seeded by large-scale asymmetries from the early explosion phase, these develop into the finger-like structures characteristic of the RayleighTaylor instability, and there is also considerable growth of small-scale plumes. Image reproduced with permission from Wongwathanarat et al. (2015), copyright by ESO

<!-- image -->

gradient is established and matter is decelerated behind the shock (Chevalier 1976). Generally, the shock velocity v sh also decreases 17 as the mass M ej of the shock ejecta grows; it roughly scales as v sh / ð E exp = M ej Þ 1 = 2 . The shock can, however, transiently accelerate when it encounters density gradients steeper than q / r /C0 3 at shell interfaces (Sedov 1959). Both effects can be captured by the formula of Matzner and McKee (1999),

/C18

/C19

/C18

/C19

<!-- formula-not-decoded -->

The post-shock pressure and density profiles adjust to variations in shock velocity to establish something of a ''quasi-hydrostatic'' stratification behind the shock with an effective gravity that is directed outward . However, once the post-shock velocities become supersonic, the post-shock pressure profile can no longer globablly adjust to changing shock velocities, and reverse shocks are formed. A first reverse shock forms typically forms at a few 1000 km as the developing neutrino-heated wind crashes into more slowly moving ejecta. Later on, further reverse shocks emerge after the shock encounters various shell interfaces. Their strength depends on the

17 Note that deceleration of the post-shock matter and deceleration of the shock do not always coincide, though they are closely related phenomena. One can have \_ vr \ 0 and \_ v sh [ 0.

<!-- image -->

density jump at the shell interface. In hydrogen-rich progenitors, the reverse shock from the H/He interface is particularly strong (especially in red supergiant) and therefore sometimes referred to simply as the reverse shock.

Rayleigh-Taylor instability The non-monotonic variations in v sh result in monmonotonic post-shock entropy and density profiles, and some layers become Rayleigh-Taylor unstable (Chevalier 1976; Mu ¨ller et al. 1991; Fryxell et al. 1991). In the relevant highly compressible regime, the growth rate for the Rayleigh-Taylor instability from a local stability analysis is given by Bandiera (1984), Benz and Thielemann (1990) and Mu ¨ller et al. (1991)

s

ffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffi

/C18

ffi

/C19

ffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffiffi

/C18

ffi

/C19

<!-- formula-not-decoded -->

where g eff ¼ q /C0 1 o P = o r is the effective gravity. The second form elucidates that stability is determined by the sub- or superadiabaticity of the density gradient as for buoyancy-driven convection. In the relevant radiation-dominated regime, composition gradients have a minor impact on stability, and the entropy gradient is the deciding factor. However, since the effective gravity points outwards, positive entropy gradients are now unstable. Such positive entropy gradients arise when the shock accelerates at shell interfaces. One should note that the actual growth rate of perturbations depends on their length scale (Zhou 2017), and Eq. (68) roughly applies to the fastest growing modes with a wavelength comparable to the width of the unstable region. Since the unstable regions tend to be narrow, Rayleigh-Taylor mixing tends to produce smaller, clumpy structures, but large-scale asymmetries can also grow considerably for sufficiently strong seed perturbations.

Intriguingly, it has been suggested that the overall effect of Rayleigh-Taylor mixing can roughly be captured in 1D by an appropriate turbulence model (Duffell 2016; Paxton et al. 2018). The key idea here is to incorporate the proper growth rate, saturation behavior, and a velocity-dependent mixing length (Duffell 2016). While a first comparison of this model with 3D results from Wongwathanarat et al. (2015) proved encouraging (Paxton et al. 2018), a few caveats remain. A detailed analysis of Rayleigh-Taylor mixing in a 3D model of a stripped-envelope progenitor by Mu ¨ller et al. (2017a) unearthed some basis for a phenomenological 1D description of Rayleigh-Taylor mixing (most notably buoyancy-drag balance in the non-linear stage) and suggested some improvements to the model of Paxton et al. (2018), but cast doubt on the use of local gradient to estimate the density and composition contrasts of the Rayleigh-Taylor plumes. In particular, the RayleighTaylor instability sometimes produces partial inversions of the initial composition profiles, which cannot be modeled by diffusive mixing in 1D.

In addition to the Rayleigh-Taylor instability, the Richtmyer-Meshkov instability (see Richtmyer 1960; Zhou 2017, for details of the instability mechanism) can develop because the shock is generally aspherical and hits the shell interfaces obliquely. The literature on mixing instabilities in supernovae is extensive, and we can only provide a very condensed summary of extant numerical studies. We will exclusively focus on the optically thick phase of the explosion and not consider the remnant phase.

<!-- image -->

s

3

Simulations of mixing in SN 1987A After a few earlier numerical experiments, significant interest in mixing instabilities was prompted by observations of SN 1987A that pointed to strong early mixing of nickel (see Arnett et al. 1989b; McCray 1993, and references therein). Two-dimensional simulations of mixing instabilities (Arnett et al. 1989a; Mu ¨ller et al. 1991; Fryxell et al. 1991; Hachisu et al. 1990; Benz and Thielemann 1990; Herant and Benz 1992) in the wake of SN 1987A took a first step towards explaining the observed mixing. The typical picture revealed by these models is that of a strong Rayleigh-Taylor instability at the H/He-interface with linear growth factors of thousands (Mu ¨ller et al. 1991), Some models (Mu ¨ller et al. 1991; Fryxell et al. 1991) also showed a second strongly unstable region at the He/C-interface that eventually merges with the mixed region further outside. Mixing was dominated by many small-scale plumes in these first-generation simulations. However, the maximum velocities of nickel plumes still fell short by about a factor of two compared to the observed velocities of up to /C24 4000 km s /C0 1 .

Many subsequent studies have investigated stronger, large-scale initial seed perturbations as a possible explanations for the strong mixing in SN 1987A and other observed Type II supernovae. Such seed perturbations are naturally expected in the neutrino-driven paradigm from the SASI and low-' convective modes, and in magnetorotational explosions with veritable jets. Most simulations have explored the effect of large-scale seed perturbations by specifying them ad hoc (e.g., Nagataki et al. 1998; Nagataki 2000; Hungerford et al. 2003; Couch et al. 2009; Ono et al. 2013; Mao et al. 2015; Ellinger et al. 2012). One must therefore be cautious in drawing conclusions on the role of ''jets'' in explaining the observed mixing. In fact, simulations with artificially injected jets rather serve to rule out kinetically-dominated jets in Type IIP supernovae based on early spectropolarimetry, though thermally-dominated jets are not excluded in principle (Couch et al. 2009).

In their 2D AMR simulations, Kifonidis et al. (2000, 2003, 2006) followed a more consistent approach by starting from light-bulb simulations of neutrino-driven explosions. While the seed asymmetries from the early explosion phase initially led to high nickel clump velocities in the Type IIP model of Kifonidis et al. (2000, 2003), the final velocities were still too small because the clumps were caught behind the reverse shock and underwent fast deceleration by supersonic drag after crossing it. Kifonidis et al. (2000) found that this can be avoided with a more slowly developing and more aspherical explosion. In this case, they found less clump deceleration because the clumps make it beyond the H/He interface before the reverse shock develops, and also found strong downward mixing of hydrogen with the help of the Richtmyer-Meshkov instability.

A very convincing picture of mixing in SN 1987A has emerged since the advent of 3D simulations. Simulations of single-mode perturbations by Kane et al. (2000) already suggested a faster growth of the Rayleigh-Taylor instability in 3D. A first 3D simulation of mixing in SN 1987A based on a 3D explosion model using gray

<!-- image -->

neutrino transport was conducted by Hammer et al. (2010), who were able to obtain realistic mixing of nickel, hydrogen, and other elements even without the need for strong initial shock deformation and a strong Richtmyer-Meshkov instability. Hammer et al. (2010) explain the reduced deceleration of plumes as a result of a more favorable volume-to-surface ratio of the clumps in 3D compared to 2D, where the clumps are actually toroidal. Stronger mixing in 3D was also confirmed in a study not specifically focused on SN 1987A (Joggerst et al. 2010b). The attainable nickel clump velocities are, however, quite sensitive to the progenitor structure (Wongwathanarat et al. 2015). Interestingly, the origin of the fastest and biggest clumps in Hammer et al. (2010) and in some of the other subsequent 3D simulations could be traced back to the most prominent convective bubbles that formed around shock revival, i.e., the late-time instabilities still contain traces of the early asymmetries imprinted by the neutrino-driven engine. As a next step towards model validation, synthetic light curves based on the 3D models of Wongwathanarat et al. (2015) were computed by Utrobin et al. (2015), and the results are encouraging. While the fit to the observed light curve is still not perfect, the discrepancies likely indicate uncertainties in the progenitor structure and the precise initial conditions after shock revival and not a problem of the neutrino-driven explosion scenario for SN 1987A.

Stripped-envelope supernovae Mixing instabilities are also highly relevant in the context of stripped-envelope supernovae. Due to the lack of a H envelope (or a small mass of the H envelope in the case of Type IIb events), the early asymmetries are not shredded as strongly by Rayleigh-Taylor mixing as in Type IIP supernovae, so that spectroscopy and spectropolarimetry offer a more direct glimpse on global asymmetries generated by the engine (see Wang and Wheeler 2008 for observational diagnostics of asymmetries). Moreover, the presence or absence of He lines in Ib/c supernovae is sensitive to the mixing of nickel (Dessart et al. 2012, 2015; Hachinger et al. 2012), and so is the detailed shape of the light curve (Shigeyama and Nomoto 1990; Yoon et al. 2019).

Two-dimensional simulations of Rayleigh-Taylor mixing in Ib/c supernovae were first conducted by Hachisu et al. (1991, 1994). These simulations were based on helium star models, which are a viable approximation for progenitors that lost their envelope due to Case B/Case C mass transfer, but used artificially triggered explosions. Hachisu et al. (1991, 1994) found indications of stronger mixing in less massive helium stars. Baron (1992) interpreted this as pointing towards an association of Ib and Ic supernovae with lowand high-mass helium stars, respectively. Kifonidis et al. (2000, 2003) triggered the explosion somewhat more realistically using a light-bulb scheme, but constructed their Ib supernova progenitor by artificially removing the hydrogen envelope at collapse (implying an inconsistent envelope structure); their finding on the mixing of nickel were qualitatively similar to Hachisu et al. (1991, 1994).

Wongwathanarat et al. (2017) took an ambitious step towards comparing stripped-envelope models with observations by performing a 3D simulation of a neutrino-driven explosion that matches the global asymmetries in the distribution of 44 Ti and 56 Ni and the neutron star kick in Cas A to an astonishing degree. The

<!-- image -->

required progenitor for a Type IIb (i.e., partially stripped) supernova was again constructed by manually removing part of the envelope at collapse, but in terms of simulation fidelity this is likely less of an issue than the fact that the neutrino-driven explosion was still tuned to match the desired explosion energy.

A first 3D simulation of mixing in an ultra-stripped progenitor starting from a self-consistent explosion model was conducted by Mu ¨ller et al. (2018) with a view to observations of fast and faint Ic supernovae (Drout et al. 2013; De et al. 2018). The model showed mixing of substantial amounts of nickel in a few narrow dense plumes out to about half way through the He envelope. These findings are, however, difficult to extrapolate to other, less extreme stripped-envelope supernova progenitors. A more thorough exploration of mixing in Ib/c supernovae and a quantitative comparison of 3D models of mixing instabilities with the spectropolarimtery of observed Ib/c events is called for.

Mixing and fallback Mixing instabilities have also been studied as a possible explanation of abundances in extremely metal-poor stars. Umeda and Nomoto (2003) suggested that the high [C/Fe] in some of these stars can be explained by invoking a combination of Rayleigh-Taylor mixing and fallback in the supernovae that supposedly contributed to their initial composition. 18 Joggerst et al. (2009) conducted 2D simulations of this scenario using the FLASH code. Their simulations indeed showed enhanced fallback in lowand zero-metallicity progenitors, and hence a possible mechanism for low iron-group yields in metal-poor environments, but Rayleigh-Taylor mixing was not sufficient for ejecting the required amount of iron-group and intermediate-mass elements to match observed abundances. In a follow-up study that surveyed a broader range of rotating progenitors with two different metallicities ( Z ¼ 0 and Z ¼ 10 /C0 4 Z /C12 ) using the CASTRO code, Joggerst et al. (2010a) were able to find better matches of the supernova yields to abundance patterns from ultra-metal poor stars. A similar study was conducted by Chen et al. (2017) to explain the abundances in the most iron-poor star to date (SMSS J031300.36-670839.3, Keller et al. 2014) by fallback in an explosion with modest energy. However, all of these simulations were restricted to 2D and imposed seed perturbations by hand in spherically symmetric models. A first 3D simulation of a fallback supernova from collapse to shock revival by the neutrino-driven mechanism, through black hole formation, and on to shock breakout was performed by Chan et al. (2018). In their model, fallback proceeds in a qualitatively different manner than in previous studies; the iron-group material is accreted early, the postshock flow involves global asymmetries during the first tens of seconds (which could potentially generate substantial black hole kicks and spins), but no mixing instabilities occur later on. While the work of Chan et al. (2018) has demonstrated the feasibility of a forward-modelling approach to fallback supernovae, they explored only a single progenitor, and a broader investigation is necessary to understand the phenomenology of fallback in three dimensions.

18 Jet-driven explosions could provide an alternative mechanism to explain the observed abundance patterns (e.g., Maeda and Nomoto 2003; Nomoto et al. 2006).

<!-- image -->

Acknowledgements I wish to thank S. Bruenn, A. Burrows, C. Collins, J. Guilet, T. Foglizzo, H.Th. Janka, K. Kotake, E. Lentz, A. Mezzacappa, E. Mu ¨ller, J. Powell, A. Skinner, T. Takiwaki, and A. Wongwathanarat for their kind permission to reproduce figures from their papers. I am grateful to A. Heger for providing 1D progenitor models at various evolutionary stages. This work was supported by the Australian Research Council (ARC) through Future Fellowship FT160100035. The author acknowledges support from the ARC Centre of Excellence for Gravitational Wave Discovery (OzGrav) as associate investigator (Grant No. CE170100004). This research was undertaken with the assistance of resources obtained via NCMAS (Grant No. fh6) and ASTAC from the National Computational Infrastructure (NCI), which is supported by the Australian Government and was supported by resources provided by the Pawsey Supercomputing Centre with funding from the Australian Government and the Government of Western Australia.

Open Access This article is licensed under a Creative Commons Attribution 4.0 International License, which permits use, sharing, adaptation, distribution and reproduction in any medium or format, as long as you give appropriate credit to the original author(s) and the source, provide a link to the Creative Commons licence, and indicate if changes were made. The images or other third party material in this article are included in the article's Creative Commons licence, unless indicated otherwise in a credit line to the material. If material is not included in the article's Creative Commons licence and your intended use is not permitted by statutory regulation or exceeds the permitted use, you will need to obtain permission directly from the copyright holder. To view a copy of this licence, visit http:// creativecommons.org/licenses/by/4.0/.

## References

- Abdikamalov E, Foglizzo T (2020) Acoustic wave generation in collapsing massive stars with convective shells. Mon Not R Astron Soc 493(3):3496-3512. https://doi.org/10.1093/mnras/staa533. arXiv: 1907.06966
- Abdikamalov E, Ott CD, Radice D, Roberts LF, Haas R, Reisswig C, Mo ¨sta P, Klion H, Schnetter E (2015) Neutrino-driven turbulent convection and standing accretion shock instability in threedimensional core-collapse supernovae. Astrophys J 808:70. https://doi.org/10.1088/0004-637X/808/ 1/70
- Abdikamalov E, Zhaksylykov A, Radice D, Berdibek S (2016) Shock-turbulence interaction in corecollapse supernovae. Mon Not R Astron Soc 461:3864-3876. https://doi.org/10.1093/mnras/ stw1604
- Abdikamalov E, Huete C, Nussupbekov A, Berdibek S (2018) Turbulence generation by shock-acousticwave interaction in core-collapse supernovae. arXiv:1805.03957
- Akiyama S, Wheeler JC, Meier DL, Meier Lichtenstadt DL, Lichtenstadt I (2003) The magnetorotational instability in core-collapse supernova explosions. Astrophys J 584:954-970. https://doi.org/10.1086/ 344135
- Almansto ¨tter M, Melson T, Janka HT, Mu ¨ller E (2018) Parallelized solution method of the threedimensional gravitational potential on the Yin-Yang grid. Astrophys J 863(2):142. https://doi.org/ 10.3847/1538-4357/aad33a. arXiv:1806.04593
- Almgren AS, Bell JB, Rendleman CA, Zingale M (2006) Low mach number modeling of type Ia supernovae. I. Hydrodynamics. Astrophys J 637:922-936. https://doi.org/10.1086/498426
- Andrassy R, Herwig F, Woodward P, Ritter C (2020) 3D hydrodynamic simulations of C ingestion into a convective O shell. Mon Not R Astron Soc 491:972-992. https://doi.org/10.1093/mnras/stz2952. arXiv:1808.04014
- Andresen H, Mu ¨ller B, Mu ¨ller E, Janka HT (2017) Gravitational wave signals from 3D neutrino hydrodynamics simulations of core-collapse supernovae. Mon Not R Astron Soc 468(2):2032-2051. https://doi.org/10.1093/mnras/stx618
- Antoniadis J, Tauris TM, Ozel F, Barr E, Champion DJ, Freire PCC (2016) The millisecond pulsar mass distribution: evidence for bimodality and constraints on the maximum neutron star mass. arXiv: 1605.01665

<!-- image -->

- Ardeljan NV, Bisnovatyi-Kogan GS, Moiseenko SG (2005) Magnetorotational supernovae. Mon Not R Astron Soc 359(1):333-344. https://doi.org/10.1111/j.1365-2966.2005.08888.x
- Arnett D (1994) Oxgen-burning hydrodynamics. 1: steady shell burning. Astrophys J 427:932-946. https://doi.org/10.1086/174199
- Arnett D (1996) Supernovae and nucleosynthesis: an investigation of the history of matter from the big bang to the present. Princeton University Press, Princeton
- Arnett WD, Meakin C (2010) Turbulent mixing in stars: theoretical hurdles. In: Cunha K, Spite M, Barbuy B (eds) Chemical abundances in the universe: connecting first stars to planets, IAU symposium, vol 265, pp 106-110. https://doi.org/10.1017/S174392131000030X
- Arnett WD, Meakin C (2011a) Toward realistic progenitors of core-collapse supernovae. Astrophys J 733:78. https://doi.org/10.1088/0004-637X/733/2/78. arXiv:1101.5646
- Arnett WD, Meakin C (2011b) Turbulent cells in stars: fluctuations in kinetic energy and luminosity. Astrophys J 741:33. https://doi.org/10.1088/0004-637X/741/1/33
- Arnett D, Fryxell B, Mueller E (1989a) Instabilities and nonradial motion in SN 1987A. Astrophys J 341:L63. https://doi.org/10.1086/185458
- Arnett WD, Bahcall JN, Kirshner RP, Woosley SE (1989b) Supernova 1987A. Annu Rev Astron Astrophys 27:629-700. https://doi.org/10.1146/annurev.aa.27.090189.003213
- Arnett D, Meakin C, Young PA (2009) Turbulent convection in stellar interiors. II. The velocity field. Astrophys J 690:1715-1729. https://doi.org/10.1088/0004-637X/690/2/1715
- Arnett WD, Meakin C, Viallet M, Campbell SW, Lattanzio JC, Moca ´k M (2015) Beyond mixing-length theory: a step Toward 321D. Astrophys J 809:30. https://doi.org/10.1088/0004-637X/809/1/30
- Asida SM, Arnett D (2000) Further adventures: oxygen burning in a convective shell. Astrophys J 545:435-443. https://doi.org/10.1086/317774. arXiv:astro-ph/0006451
- Baade W, Zwicky F (1934) On super-novae. Proc Natl Acad Sci USA 20:254-259. https://doi.org/10. 1073/pnas.20.5.254
- Balbus SA, Hawley JF (1991) A powerful local shear instability in weakly magnetized disks. I. Linear analysis. II. Nonlinear evolution. Astrophys J 376:214-233. https://doi.org/10.1086/170270
- Balsara DS (2017) Higher-order accurate space-time schemes for computational astrophysics-Part I: finite volume methods. Living Rev Comput Astrophys 3:2. https://doi.org/10.1007/s41115-0170002-8
- Bandiera R (1984) Convective supernovae. Astron Astrophys 139:368-374
- Baron E (1992) Progenitor masses of type Ib/c supernovae. Mon Not R Astron Soc 255:267-268. https:// doi.org/10.1093/mnras/255.2.267
- Baumgarte TW, Shapiro SL (2010) Numerical relativity: solving Einstein's equations on the computer. Cambridge University Press, New York
- Baumgarte TW, Montero PJ, Cordero-Carrio ´n I, Mu ¨ller E (2013) Numerical relativity in spherical polar coordinates: evolution calculations with the BSSN formulation. Phys Rev D 87(4):044026. https:// doi.org/10.1103/PhysRevD.87.044026
- Baumgarte TW, Montero PJ, Mu ¨ller E (2015) Numerical relativity in spherical polar coordinates: offcenter simulations. Phys Rev D 91(6):064035. https://doi.org/10.1103/PhysRevD.91.064035
- Bazan G, Arnett D (1994) Convection, nucleosynthesis, and core collapse. Astrophys J 433:L41-L43. https://doi.org/10.1086/187543
- Baza ´n G, Arnett D (1997) Large nuclear networks in presupernova models. Nucl Phys A 621:607-610. https://doi.org/10.1016/S0375-9474(97)00313-8
- Bazan G, Arnett D (1998) Two-dimensional hydrodynamics of pre-core collapse: oxygen shell burning. Astrophys J 496:316. https://doi.org/10.1086/305346. arXiv:astro-ph/9702239
- Benz W, Thielemann FK (1990) Convective instabilities in SN 1987A. Astrophys J 348:L17-L20. https:// doi.org/10.1086/185620
- Berger MJ, Colella P (1989) Local adaptive mesh refinement for shock hydrodynamics. J Comput Phys 82:64-84. https://doi.org/10.1016/0021-9991(89)90035-1
- Bethe HA (1990) Supernova mechanisms. Rev Mod Phys 62:801-866. https://doi.org/10.1103/ RevModPhys.62.801
- Bethe HA, Brown GE, Cooperstein J (1987) Convection in supernova theory. Astrophys J 322:201. https://doi.org/10.1086/165715
- Biermann L (1932) Untersuchungen u ¨ber den inneren Aufbau der Sterne. IV. Konvektionszonen im Innern der Sterne. Z Astrophys 5:117 (Vero ¨ffentlichungen der Universita ¨ts-Sternwarte Go ¨ttingen, Nr. 27)

<!-- image -->

- Blondin JM, Lufkin EA (1993) The piecewise-parabolic method in curvilinear coordinates. Astrophys J Suppl 88:589. https://doi.org/10.1086/191834
- Blondin JM, Mezzacappa A (2006) The spherical accretion shock instability in the linear regime. Astrophys J 642:401-409. https://doi.org/10.1086/500817. arXiv:astro-ph/0507181
- Blondin JM, Mezzacappa A (2007) Pulsar spins from an instability in the accretion shock of supernovae. Nature 445:58-60. https://doi.org/10.1038/nature05428. arXiv:astro-ph/0611680
- Blondin JM, Mezzacappa A, DeMarino C (2003) Stability of standing accretion shocks, with an eye toward core-collapse supernovae. Astrophys J 584:971-980. https://doi.org/10.1086/345812. arXiv: astro-ph/0210634
- Blondin JM, Gipson E, Harris S, Mezzacappa A (2017) The standing accretion shock instability: enhanced growth in rotating progenitors. Astrophys J 835(2):170. https://doi.org/10.3847/15384357/835/2/170
- Bludman SA, van Riper KA (1978) Diffusion approximation to neutrino transport in dense matter. Astrophys J 224:631-642. https://doi.org/10.1086/156412
- Bodansky D, Clayton DD, Fowler WA (1968) Nuclear quasi-equilibrium during silicon burning. Astrophys J Suppl 16:299. https://doi.org/10.1086/190176
- Bo ¨ hm-Vitense E (1958) U ¨ ber die Wasserstoffkonvektionszone in Sternen verschiedener Effektivtemperaturen und Leuchtkra ¨fte. Z Astrophys 46:108
- Bollig R, Janka HT, Lohs A, Martı ´nez-Pinedo G, Horowitz CJ, Melson T (2017) Muon creation in supernova matter facilitates neutrino-driven explosions. Phys Rev Lett 119(24):242702. https://doi. org/10.1103/PhysRevLett.119.242702
- Bonazzola S, Gourgoulhon E, Grandcle ´ment P, Novak J (2004) Constrained scheme for the Einstein equations based on the Dirac gauge and spherical coordinates. Phys Rev D 70(10):104007. https:// doi.org/10.1103/PhysRevD.70.104007. arXiv:gr-qc/0307082
- Boyd JP (2001) Chebyshev and Fourier spectral methods, 2nd edn. Dover, Mineola
- Bray JC, Eldridge JJ (2016) Neutron star kicks and their relationship to supernovae ejecta mass. Mon Not R Astron Soc 461:3747-3759. https://doi.org/10.1093/mnras/stw1275
- Bruenn SW, Dineva T (1996) The role of doubly diffusive instabilities in the core-collapse supernova mechanism. Astrophys J 458:L71-L74. https://doi.org/10.1086/309921
- Bruenn SW, Mezzacappa A, Dineva T (1995) Dynamic and diffusive instabilities in core collapse supernovae. Phys Rep 256:69-94. https://doi.org/10.1016/0370-1573(94)00102-9
- Bruenn SW, Raley EA, Mezzacappa A (2004) Fluid stability below the neutrinospheres of supernova progenitors and the dominant role of lepto-entropy fingers. arXiv:astro-ph/0404099
- Bruenn SW, Mezzacappa A, Hix WR, Lentz EJ, Messer OEB, Lingerfelt EJ, Blondin JM, Endeve E, Marronetti P, Yakunin KN (2013) Axisymmetric ab initio core-collapse supernova simulations of 12-25 M /C12 stars. Astrophys J 767(1):L6. https://doi.org/10.1088/2041-8205/767/1/L6
- Bruenn SW, Lentz EJ, Hix WR, Mezzacappa A, Harris JA, Messer OEB, Endeve E, Blondin JM, Chertkow MA, Lingerfelt EJ, Marronetti P, Yakunin KN (2016) The development of explosions in axisymmetric ab initio core-collapse supernova simulations of 12-25 M /C12 stars. Astrophys J 818:123. https://doi.org/10.3847/0004-637X/818/2/123
- Bruenn SW, Blondin JM, Hix WR, Lentz EJ, Messer OEB, Mezzacappa A, Endeve E, Harris JA, Marronetti P, Budiardja RD, Chertkow MA, Lee CT (2018) Chimera: a massively parallel code for core-collapse supernova simulation. arXiv:1809.05608
- Bru ¨ggen M, Hillebrandt W (2001) Three-dimensional simulations of shear instabilities in magnetized flows. Mon Not R Astron Soc 323:56-66. https://doi.org/10.1046/j.1365-8711.2001.04206.x
- Bugli M, Guilet J, Obergaulinger M, Cerda ´-Dura ´n P, Aloy MA (2020) The impact of non-dipolar magnetic fields in core-collapse supernovae. Mon Not R Astron Soc 492(1):58-71. https://doi.org/ 10.1093/mnras/stz34. arXiv:1909.02824
- Buras R, Janka HT, Rampp M, Kifonidis K (2006a) Two-dimensional hydrodynamic core-collapse supernova simulations with spectral neutrino transport. II. Models for different progenitor stars. Astron Astrophys 457:281-308. https://doi.org/10.1051/0004-6361:20054654. arXiv:astro-ph/ 0512189
- Buras R, Rampp M, Janka HT, Kifonidis K (2006b) Two-dimensional hydrodynamic core-collapse supernova simulations with spectral neutrino transport. I. Numerical method and results for a 15 M /C12 star. Astron Astrophys 447:1049-1092. https://doi.org/10.1051/0004-6361:20053783. arXiv:astroph/0507135
- Burrows A (1987) Convection and the mechanism of type II supernovae. Astrophys J 318:L57-L61. https://doi.org/10.1086/184937

<!-- image -->

- Burrows A (2013) Colloquium: perspectives on core-collapse supernova theory. Rev Mod Phys 85:245-261. https://doi.org/10.1103/RevModPhys.85.245. arXiv:1210.4921
- Burrows A, Goshy J (1993) A theory of supernova explosions. Astrophys J 416:L75 ? . https://doi.org/10. 1086/187074
- Burrows A, Hayes J (1996) Pulsar recoil and gravitational radiation due to asymmetrical stellar collapse and explosion. Phys Rev Lett 76:352-355. https://doi.org/10.1103/PhysRevLett.76.352. arXiv:astroph/9511106
- Burrows A, Lattimer JM (1985) The prompt mechanism of type II supernovae. Astrophys J 299:L19-L22. https://doi.org/10.1086/184572
- Burrows A, Lattimer JM (1988) Convection, type II supernovae, and the early evolution of neutron stars. Phys Rep 163:51-62. https://doi.org/10.1016/0370-1573(88)90035-X
- Burrows A, Hayes J, Fryxell BA (1995) On the nature of core-collapse supernova explosions. Astrophys J 450:830-850. https://doi.org/10.1086/176188. arXiv:astro-ph/9506061
- Burrows A, Dessart L, Livne E, Ott CD, Murphy J (2007a) Simulations of magnetically driven supernova and hypernova explosions in the context of rapid rotation. Astrophys J 664:416-434. https://doi.org/ 10.1086/519161. arXiv:astro-ph/0702539
- Burrows A, Livne E, Dessart L, Ott CD, Murphy J (2007b) Features of the acoustic mechanism of corecollapse supernova explosions. Astrophys J 655:416-433. https://doi.org/10.1086/509773. arXiv: astro-ph/0610175
- Burrows A, Vartanyan D, Dolence JC, Skinner MA, Radice D (2018) Crucial physical dependencies of the core-collapse supernova mechanism. Space Sci Rev 214:33. https://doi.org/10.1007/s11214-0170450-9
- Burrows A, Radice D, Vartanyan D (2019) Three-dimensional supernova explosion simulations of 9-, 10-, 11-, 12-, and 13M /C12 stars. Mon Not R Astron Soc 485(3):3153-3168. https://doi.org/10.1093/ mnras/stz543
- Burrows A, Radice D, Vartanyan D, Nagakura H, Skinner MA, Dolence JC (2020) The overarching framework of core-collapse supernova explosions as revealed by 3D FORNAX simulations. Mon Not R Astron Soc 491(2):2715-2735. https://doi.org/10.1093/mnras/stz3223
- Calhoun DA, Helzel C, Leveque RJ (2008) Logically rectangular grids and finite volume methods for pdes in circular and spherical domains. SIAM Rev 50(4):723-752. https://doi.org/10.1137/ 060664094
- Cardall CY, Budiardja RD (2015) Stochasticity and efficiency in simplified models of core-collapse supernova explosions. Astrophys J 813:L6. https://doi.org/10.1088/2041-8205/813/1/L6
- Cerda ´-Dura ´n P (2009) General relativistic hydrodynamics beyond bouncing polytropes. https://indico.nbi. ku.dk/event/50/contributions/2696/
- Cerda ´-Dura ´n P, Faye G, Dimmelmeier H, Font JA, Iba ´n ˜ez JM, Mu ¨ller E, Scha ¨fer G (2005) CFC ? : improved dynamics and gravitational waveforms from relativistic core collapse simulations. Astron Astrophys 439(3):1033-1055. https://doi.org/10.1051/0004-6361:20042602
- Chan C, Mu ¨ller B, Heger A, Pakmor R, Springel V (2018) Black hole formation and fallback during the supernova explosion of a 40 M /C12 star. Astrophys J 852:L19. https://doi.org/10.3847/2041-8213/ aaa28c
- Chan C, Mu ¨ller B, Heger A (2020) The impact of fallback on the compact remnants and chemical yields of core-collapse supernovae. arXiv e-prints arXiv:2003.04320
- Chandrasekhar S (1961) Hydrodynamic and hydromagnetic stability. Clarendon, Oxford
- Chatzopoulos E, Graziani C, Couch SM (2014) Characterizing the convective velocity fields in massive stars. Astrophys J 795:92. https://doi.org/10.1088/0004-637X/795/1/92
- Chatzopoulos E, Couch SM, Arnett WD, Timmes FX (2016) Convective properties of rotating twodimensional core-collapse supernova progenitors. Astrophys J 822:61. https://doi.org/10.3847/0004637X/822/2/61
- Chen KJ, Heger A, Almgren AS (2013) Numerical approaches for multidimensional simulations of stellar explosions. Astron Comput 3:70-78. https://doi.org/10.1016/j.ascom.2014.01.001
- Chen KJ, Heger A, Whalen DJ, Moriya TJ, Bromm V, Woosley SE (2017) Low-energy population III supernovae and the origin of extremely metal-poor stars. Mon Not R Astron Soc 467(4):4731-4738. https://doi.org/10.1093/mnras/stx470
- Chevalier RA (1976) The hydrodynamics of type II supernovae. Astrophys J 207:872-887. https://doi. org/10.1086/154557
- Clayton DD (1968) Principles of stellar evolution and nucleosynthesis. McGraw-Hill, New York

<!-- image -->

- Clercx HJH, van Heijst GJF (2017) Dissipation of coherent structures in confined two-dimensional turbulence. Phys Fluids 29(11):111103. https://doi.org/10.1063/1.4993488
- Colella P, Glaz HM (1985) Efficient solution algorithms for the Riemann problem for real gases. J Comput Phys 59:264-289
- Colella P, Sekora MD (2008) A limiter for PPM that preserves accuracy at smooth extrema. J Comput Phys 227:7069-7076. https://doi.org/10.1016/j.jcp.2008.03.034
- Colella P, Woodward PR (1984) The piecewise parabolic method (PPM) for gas-dynamical simulations. J Comput Phys 54:174-201
- Colella P, Majda A, Roytburd V (1986) Theoretical and numerical structure for reacting shock waves. SIAM J Sci Stat Comput 7(4):1059-1080. https://doi.org/10.1137/0907073
- Colgate SA, White RH (1966) The hydrodynamic behavior of supernovae explosions. Astrophys J 143:626-681. https://doi.org/10.1086/148549
- Colgate SA, Grasberger WH, White RH (1961) The dynamics of supernova explosions. Astron J 70:280. https://doi.org/10.1086/108573
- Collins C, Mu ¨ller B, Heger A (2018) Properties of convective oxygen and silicon burning shells in supernova progenitors. Mon Not R Astron Soc 473:1695-1704. https://doi.org/10.1093/mnras/ stx2470
- Cordero-Carrio ´n I, Cerda ´-Dura ´n P, Dimmelmeier H, Jaramillo JL, Novak J, Gourgoulhon E (2009) Improved constrained scheme for the Einstein equations: an approach to the uniqueness issue. Phys Rev D 79(2):024017. https://doi.org/10.1103/PhysRevD.79.024017. arXiv:0809.2325
- Cordero-Carrio ´n I, Cerda ´-Dura ´n P, Iba ´n ˜ez JM (2012) Gravitational waves in dynamical spacetimes with matter content in the fully constrained formulation. Phys Rev D 85(4):044023. https://doi.org/10. 1103/PhysRevD.85.044023. arXiv:1108.0571
- Co ˆ te ´ B, Jones S, Herwig F, Pignatari M (2020) Chromium nucleosynthesis and silicon-carbon shell mergers in massive stars. Astrophys J 892(1):57. https://doi.org/10.3847/1538-4357/ab77ac. arXiv: 1906.07218
- Couch SM (2013) On the impact of three dimensions in simulations of neutrino-driven core-collapse supernova explosions. Astrophys J 775:35. https://doi.org/10.1088/0004-637X/775/1/35
- Couch SM (2017) The mechanism(s) of core-collapse supernovae. Philos Trans R Soc London Ser A 375(2105):20160271. https://doi.org/10.1098/rsta.2016.0271
- Couch SM, O'Connor EP (2014) High-resolution three-dimensional simulations of core-collapse supernovae in multiple progenitors. Astrophys J 785:123. https://doi.org/10.1088/0004-637X/785/2/ 123
- Couch SM, Ott CD (2013) Revival of the stalled core-collapse supernova shock triggered by precollapse asphericity in the progenitor star. Astrophys J 778:L7. https://doi.org/10.1088/2041-8205/778/1/L7. arXiv:1309.2632
- Couch SM, Ott CD (2015) The role of turbulence in neutrino-driven core-collapse supernova explosions. Astrophys J 799:5. https://doi.org/10.1088/0004-637X/799/1/5
- Couch SM, Wheeler JC, Milosavljevic ´ M (2009) Aspherical core-collapse supernovae in red supergiants powered by nonrelativistic jets. Astrophys J 696(1):953-970. https://doi.org/10.1088/0004-637X/ 696/1/953
- Couch SM, Graziani C, Flocke N (2013) An improved multipole approximation for self-gravity and its importance for core-collapse supernova simulations. Astrophys J 778:181. https://doi.org/10.1088/ 0004-637X/778/2/181
- Couch SM, Chatzopoulos E, Arnett WD, Timmes FX (2015) The three-dimensional evolution to core collapse of a massive star. Astrophys J 808:L21. https://doi.org/10.1088/2041-8205/808/1/L21
- Courant R, Friedrichs K, Lewy H (1928) U ¨ ber die partiellen Differenzengleichungen der mathematischen Physik. Math Ann 100:32-74. https://doi.org/10.1007/BF01448839
- Cristini A, Meakin C, Hirschi R, Arnett D, Georgy C, Viallet M (2016) Linking 1D evolutionary to 3D hydrodynamical simulations of massive stars. Phys Scripta 91(3):034006. https://doi.org/10.1088/ 0031-8949/91/3/034006
- Cristini A, Meakin C, Hirschi R, Arnett D, Georgy C, Viallet M, Walkington I (2017) 3D hydrodynamic simulations of carbon burning in massive stars. Mon Not R Astron Soc 471(1):279-300. https://doi. org/10.1093/mnras/stx1535
- Cristini A, Hirschi R, Meakin C, Arnett D, Georgy C, Walkington I (2019) Dependence of convective boundary mixing on boundary properties and turbulence strength. Mon Not R Astron Soc 484(4):4645-4664. https://doi.org/10.1093/mnras/stz312

<!-- image -->

3

- Davis A, Jones S, Herwig F (2019) Convective boundary mixing in a post-He core burning massive star model. Mon Not R Astron Soc 484(3):3921-3934. https://doi.org/10.1093/mnras/sty3415
- De K, Kasliwal MM, Ofek EO, Moriya TJ, Burke J, Cao Y, Cenko SB, Doran GB, Duggan GE, Fender RP, Fransson C, Gal-Yam A, Horesh A, Kulkarni SR, Laher RR, Lunnan R, Manulis I, Masci F, Mazzali PA, Nugent PE, Perley DA, Petrushevska T, Piro AL, Rumsey C, Sollerman J, Sullivan M, Taddia F (2018) A hot and fast ultra-stripped supernova that likely formed a compact neutron star binary. Science 362(6411):201-206. https://doi.org/10.1126/science.aas8693
- Dearborn DSP, Lattanzio JC, Eggleton PP (2006) Three-dimensional numerical experimentation on the core helium flash of low-mass red giants. Astrophys J 639(1):405-415. https://doi.org/10.1086/ 499263
- Despre ´s B, Labourasse E (2015) Angular momentum preserving cell-centered Lagrangian and Eulerian schemes on arbitrary grids. J Comput Phys 290:28-54. https://doi.org/10.1016/j.jcp.2015.02.032
- Dessart L, Burrows A, Livne E, Ott CD (2006) Multidimensional radiation/hydrodynamic simulations of proto-neutron star convection. Astrophys J 645:534-550. https://doi.org/10.1086/504068. arXiv: astro-ph/0510229
- Dessart L, Burrows A, Livne E, Ott CD (2007) Magnetically driven explosions of rapidly rotating white dwarfs following accretion-induced collapse. Astrophys J 669:585-599. https://doi.org/10.1086/ 521701. arXiv:0705.3678
- Dessart L, Hillier DJ, Li C, Woosley S (2012) On the nature of supernovae Ib and Ic. Mon Not R Astron Soc 424:2139-2159. https://doi.org/10.1111/j.1365-2966.2012.21374.x
- Dessart L, Hillier DJ, Woosley S, Livne E, Waldman R, Yoon SC, Langer N (2015) Radiative-transfer models for supernovae IIb/Ib/Ic from binary-star progenitors. Mon Not R Astron Soc 453:2189-2213. https://doi.org/10.1093/mnras/stv1747
- Dimmelmeier H, Font JA, Mu ¨ller E (2002) Relativistic simulations of rotational core collapse I. Methods, initial models, and code tests. Astron Astrophys 388:917-935. https://doi.org/10.1051/0004-6361: 20020563
- Dimmelmeier H, Novak J, Font JA, Iba ´n ˜ez JM, Mu ¨ller E (2005) Combining spectral and shock-capturing methods: a new numerical approach for 3D relativistic core collapse simulations. Phys Rev D 71(064023):1-30. https://doi.org/10.1103/PhysRevD.71.064023
- Doherty CL, Gil-Pons P, Siess L, Lattanzio JC (2017) Super-AGB stars and their role as electron capture supernova progenitors. Publ Astron Soc Australia 34:e056. https://doi.org/10.1017/pasa.2017.52
- Dolence JC, Burrows A, Murphy JW, Nordhaus J (2013) Dimensional dependence of the hydrodynamics of core-collapse supernovae. Astrophys J 765:110. https://doi.org/10.1088/0004-637X/765/2/110. arXiv:1210.5241
- Donat R, Marquina A (1996) Capturing shock reflections: an improved flux formula. J Comput Phys 125:42-58. https://doi.org/10.1006/jcph.1996.0078
- Drout MR, Soderberg AM, Mazzali PA, Parrent JT, Margutti R, Milisavljevic D, Sanders NE, Chornock R, Foley RJ, Kirshner RP, Filippenko AV, Li W, Brown PJ, Cenko SB, Chakraborti S, Challis P, Friedman A, Ganeshalingam M, Hicken M, Jensen C, Modjaz M, Perets HB, Silverman JM, Wong DS (2013) The fast and furious decay of the peculiar type Ic supernova 2005ek. Astrophys J 774:58. https://doi.org/10.1088/0004-637X/774/1/58
- Duffell PC (2016) A One-dimensional model for Rayleigh-Taylor instability in supernova remnants. Astrophys J 821:76. https://doi.org/10.3847/0004-637X/821/2/76
- Duffell PC, MacFadyen AI (2011) TESS: a relativistic hydrodynamics code on a moving voronoi mesh. Astrophys J Suppl 197(2):15. https://doi.org/10.1088/0067-0049/197/2/15
- Duffell PC, MacFadyen AI (2013) Rayleigh-Taylor instability in a relativistic fireball on a moving computational grid. Astrophys J 775(2):87. https://doi.org/10.1088/0004-637X/775/2/87
- Duffell PC, MacFadyen AI (2015) From engine to afterglow: collapsars naturally produce top-heavy jets and early-time plateaus in gamma-ray burst afterglows. Astrophys J 806(2):205. https://doi.org/10. 1088/0004-637X/806/2/205
- Eastwood JW, Brownrigg DRK (1979) Remarks on the solution of Poisson's equation for isolated systems. J Comput Phys 32:24-38. https://doi.org/10.1016/0021-9991(79)90139-6
- Einfeldt B (1988) On Godunov-type methods for gas dynamics. SIAM J Numer Anal 25:294-318
- Ellinger CI, Young PA, Fryer CL, Rockefeller G (2012) A case study of small-scale structure formation in three-dimensional supernova simulations. Astrophys J 755(2):160. https://doi.org/10.1088/0004637X/755/2/160
- Endeve E, Cardall CY, Budiardja MA (2010) Generation of magnetic fields by the stationary accretion shock instability. Astrophys J 713:1219-1243. https://doi.org/10.1088/0004-637X/713/2/1219

<!-- image -->

- Endeve E, Cardall CY, Budiardja RD, Beck SW, Bejnood A, Toedte RJ, Mezzacappa A, Blondin JM (2012) Turbulent magnetic field amplification from spiral SASI modes: implications for corecollapse supernovae and proto-neutron star magnetization. Astrophys J 751:26. https://doi.org/10. 1088/0004-637X/751/1/26
- Epstein RI (1979) Lepton-driven convection in supernovae. Mon Not R Astron Soc 188:305-325
- Ertl T, Janka HT, Woosley SE, Sukhbold T, Ugliano M (2016) A two-parameter criterion for classifying the explodability of massive stars by the neutrino-driven mechanism. Astrophys J 818:124. https:// doi.org/10.3847/0004-637X/818/2/124
- Falkovich G, Boffetta G, Shats M, Lanotte AS (2017) Introduction to focus issue: two-dimensional turbulence. Phys Fluids 29(11):110901. https://doi.org/10.1063/1.5012997
- Faucher-Gigue `re CA, Kaspi VM (2006) Birth and evolution of isolated radio pulsars. Astrophys J 643:332-355. https://doi.org/10.1086/501516
- Ferna ´ndez R (2010) The spiral modes of the standing accretion shock instability. Astrophys J 725:1563-1580. https://doi.org/10.1088/0004-637X/725/2/1563. arXiv:1003.1730
- Ferna ´ndez R (2012) Hydrodynamics of core-collapse supernovae at the transition to explosion. I. Spherical symmetry. Astrophys J 749:142. https://doi.org/10.1088/0004-637X/749/2/142. arXiv: 1111.0665
- Ferna ´ndez R (2015) Three-dimensional simulations of SASI- and convection-dominated core-collapse supernovae. Mon Not R Astron Soc 452:2071-2086. https://doi.org/10.1093/mnras/stv1463
- Ferna ´ndez R, Thompson C (2009a) Dynamics of a spherical accretion shock with neutrino heating and alpha-particle recombination. Astrophys J 703:1464-1485. https://doi.org/10.1088/0004-637X/703/ 2/1464. arXiv:0812.4574
- Ferna ´ndez R, Thompson C (2009b) Stability of a spherical accretion shock with nuclear dissociation. Astrophys J 697:1827-1841. https://doi.org/10.1088/0004-637X/697/2/1827
- Ferna ´ndez R, Mu ¨ller B, Foglizzo T, Janka HT (2014) Characterizing SASI- and convection-dominated core-collapse supernova explosions in two dimensions. Mon Not R Astron Soc 440:2763-2780. https://doi.org/10.1093/mnras/stu408
- Fernando HJS (1991) Turbulent mixing in stratified fluids. Annu Rev Fluid Mech 23:455-493. https://doi. org/10.1146/annurev.fl.23.010191.002323
- Ferrario L, Wickramasinghe D (2006) Modelling of isolated radio pulsars and magnetars on the fossil field hypothesis. Mon Not R Astron Soc 367(3):1323-1328. https://doi.org/10.1111/j.1365-2966. 2006.10058.x
- Fischer T, Bastian NU, Blaschke D, Cierniak M, Hempel M, Kla ¨hn T, Martı ´nez-Pinedo G, Newton WG, Ro ¨ pke G, Typel S (2017) The state of matter in simulations of core-collapse supernovae-reflections and recent developments. Publ Astron Soc Australia 34:e067. https://doi.org/10.1017/pasa.2017.63
- Foglizzo T (2001) Entropic-acoustic instability of shocked Bondi accretion I. What does perturbed Bondi accretion sound like? Astron Astrophys 368:311-324. https://doi.org/10.1051/0004-6361:20000506. arXiv:astro-ph/0101056
- Foglizzo T (2002) Non-radial instabilities of isothermal Bondi accretion with a shock: vortical-acoustic cycle vs. post-shock acceleration. Astron Astrophys 392:353-368. https://doi.org/10.1051/00046361:20020912. arXiv:astro-ph/0206274
- Foglizzo T, Tagger M (2000) Entropic-acoustic instability in shocked accretion flows. Astron Astrophys 363:174-183
- Foglizzo T, Scheck L, Janka HT (2006) Neutrino-driven convection versus advection in core-collapse supernovae. Astrophys J 652:1436-1450. https://doi.org/10.1086/508443. arXiv:astro-ph/0507636
- Foglizzo T, Galletti P, Scheck L, Janka HT (2007) Instability of a stalled accretion shock: evidence for the advective-acoustic cycle. Astrophys J 654:1006-1021. https://doi.org/10.1086/509612. arXiv: astro-ph/0606640
- Fragile PC, Lindner CC, Anninos P, Salmonson JD (2009) Application of the cubed-sphere grid to tilted black hole accretion disks. Astrophys J 691(1):482-494. https://doi.org/10.1088/0004-637X/691/1/ 482
- Fransson C, Chevalier RA (1989) Late emission from supernovae: a window on stellar nucleosynthesis. Astrophys J 343:323. https://doi.org/10.1086/167707
- Freytag B, Ludwig HG, Steffen M (1996) Hydrodynamical models of stellar convection; the role of overshoot in DA white dwarfs, A-type stars, and the Sun. Astron Astrophys 313:497-516
- Fryer CL, Warren MS (2002) Modeling core-collapse supernovae in three dimensions. Astrophys J 574:L65-L68. https://doi.org/10.1086/342258

<!-- image -->

- Fryer CL, Warren MS (2004) The collapse of rotating massive stars in three dimensions. Astrophys J 601(1):391-404. https://doi.org/10.1086/380193
- Fryer CL, Rockefeller G, Warren MS (2006) SNSPH: a parallel three-dimensional smoothed particle radiation hydrodynamics code. Astrophys J 643:292-305. https://doi.org/10.1086/501493. arXiv: astro-ph/0512532
- Fryxell BA, Mu ¨ller E, Arnett D (1989) Hydrodynamics and nuclear burning. Max-Planck-Institut fu ¨r Astrophysik, Preprint, 449
- Fryxell B, Arnett D, Mueller E (1991) Instabilities and clumping in SN 1987A. I. Early evolution in two dimensions. Astrophys J 367:619-634. https://doi.org/10.1086/169657
- Fryxell BA, Olson K, Ricker P, Timmes FX, Zingale M, Lamb DQ, MacNeice P, Rosner R, Truran JW, Tufo H (2000) FLASH: an adaptive mesh hydrodynamics code for modeling astrophysical thermonuclear flashes. Astrophys J Suppl 131:273-334. https://doi.org/10.1086/317361
- Garaud P (2018) Double-diffusive convection at low Prandtl number. Annu Rev Fluid Mech 50(1):275-298. https://doi.org/10.1146/annurev-fluid-122316-045234
- Gerwin RA (1968) Stability of the interface between two fluids in relative motion. Rev Mod Phys 40(3):652-658. https://doi.org/10.1103/RevModPhys.40.652
- Gessner A, Janka HT (2018) Hydrodynamical neutron-star kicks in electron-capture supernovae and implications for the CRAB supernova. Astrophys J 865:61. https://doi.org/10.3847/1538-4357/ aadbae
- Gingold RA, Monaghan JJ (1977) Smoothed particle hydrodynamics: theory and application to nonspherical stars. Mon Not R Astron Soc 181:375-389. https://doi.org/10.1093/mnras/181.3.375
- Gizon L, Birch AC (2012) Helioseismology challenges models of solar convection. Proc Natl Acad Sci USA 109(30):11896-11897. https://doi.org/10.1073/pnas.1208875109
- Glas R, Janka HT, Melson T, Stockinger G, Just O (2019) Effects of LESA in three-dimensional supernova simulations with multidimensional and ray-by-ray-plus neutrino transport. Astrophys J 881(1):36. https://doi.org/10.3847/1538-4357/ab275c
- Glatzmaier GA (1984) Numerical simulations of stellar convective dynamos. I. The model and method. J Comput Phys 55:461-484. https://doi.org/10.1016/0021-9991(84)90033-0
- Goodwin BT, Pethick CJ (1982) Transport properties of degenerate neutrinos in dense matter. Astrophys J 253:816-838. https://doi.org/10.1086/159684
- Greenberg JM, Leroux AY (1996) A well-balanced scheme for the numerical processing of source terms in hyperbolic equations. SIAM J Num Anal 33(1):1-16
- Grefenstette BW, Harrison FA, Boggs SE, Reynolds SP, Fryer CL, Madsen KK, Wik DR, Zoglauer A, Ellinger CI, Alexander DM, An H, Barret D, Christensen FE, Craig WW, Forster K, Giommi P, Hailey CJ, Hornstrup A, Kaspi VM, Kitaguchi T, Koglin JE, Mao PH, Miyasaka H, Mori K, Perri M, Pivovaroff MJ, Puccetti S, Rana V, Stern D, Westergaard NJ, Zhang WW (2014) Asymmetries in core-collapse supernovae from maps of radioactive 44 Ti in Cassiopeia A. Nature 506:339-342. https://doi.org/10.1038/nature12997
- Gresho PM, Chan ST (1990) On the theory of semi-implicit projection methods for viscous incompressible flow and its implementation via a finite element method that also introduces a nearly consistent mass matrix. II. Implementation. Int J Num Meth Fluids 11:621-659. https://doi. org/10.1002/fld.1650110510
- Guidry MW, Billings JJ, Hix WR (2013) Explicit integration of extremely stiff reaction networks: partial equilibrium methods. Comput Sci Discovery 6(1):015003. https://doi.org/10.1088/1749-4699/6/1/ 015003
- Guilet J, Ferna ´ndez R (2014) Angular momentum redistribution by SASI spiral modes and consequences for neutron star spins. Mon Not R Astron Soc 441:2782-2798. https://doi.org/10.1093/mnras/stu718
- Guilet J, Foglizzo T (2012) On the linear growth mechanism driving the standing accretion shock instability. Mon Not R Astron Soc 421:546-560. https://doi.org/10.1111/j.1365-2966.2012.20333.x
- Guilet J, Sato J, Foglizzo T (2010) The saturation of SASI by parasitic instabilities. Astrophys J 713:1350-1362. https://doi.org/10.1088/0004-637X/713/2/1350. arXiv:0910.3953
- Guilet J, Foglizzo T, Fromang S (2011) Dynamics of an Alfve ´n surface in core collapse supernovae. Astrophys J 729(1):71. https://doi.org/10.1088/0004-637X/729/1/71
- Guilet J, Mu ¨ller E, Janka HT (2015) Neutrino viscosity and drag: impact on the magnetorotational instability in protoneutron stars. Mon Not R Astron Soc 447:3992-4003. https://doi.org/10.1093/ mnras/stu2550

<!-- image -->

- Guillard H, Murrone A (2004) On the behavior of upwind schemes in the low Mach number limit: II. Godunov type schemes. Comput Fluids 33(4):655-675. https://doi.org/10.1016/j.compfluid.2003.07. 001
- Hachinger S, Mazzali PA, Taubenberger S, Hillebrandt W, Nomoto K, Sauer DN (2012) How much H and He is 'hidden' in SNe Ib/c? I. Low-mass objects. Mon Not R Astron Soc 422:70-88. https://doi. org/10.1111/j.1365-2966.2012.20464.x
- Hachisu I, Matsuda T, Nomoto K, Shigeyama T (1990) Nonlinear growth of Rayleigh-Taylor instabilities and mixing in SN 1987A. Astrophys J 358:L57. https://doi.org/10.1086/185779
- Hachisu I, Matsuda T, Nomoto K, Shigeyama T (1991) Rayleigh-Taylor instabilities and mixing in the helium star models for type Ib/Ic supernovae. Astrophys J 368:L27. https://doi.org/10.1086/185940
- Hachisu I, Matsuda T, Nomoto K, Shigeyama T (1994) Mixing in ejecta of supernovae. II. Mixing width of 2D Rayleigh-Taylor instabilities in the helium star models for type Ib/Ic supernovae. Astron Astrophys Suppl 104:341-364
- Hammer NJ, Janka H, Mu ¨ller E (2010) Three-dimensional simulations of mixing instabilities in supernova explosions. Astrophys J 714:1371-1385. https://doi.org/10.1088/0004-637X/714/2/1371. arXiv:0908.3474
- Hanasoge SM, Duvall TL, Sreenivasan KR (2012) Anomalously weak solar convection. Proc Natl Acad Sci USA 109(30):11928-11932. https://doi.org/10.1073/pnas.1206570109
- Handy T, Plewa T, Odrzywołek A (2014) Toward Connecting core-collapse supernova theory with observations. I. Shock revival in a 15 M /C12 blue supergiant progenitor with SN 1987A energetics. Astrophys J 783:125. https://doi.org/10.1088/0004-637X/783/2/125
- Hanke F, Marek A, Mu ¨ller B, Janka HT (2012) Is strong SASI activity the key to successful neutrinodriven supernova explosions? Astrophys J 755:138. https://doi.org/10.1088/0004-637X/755/2/138. arXiv:1108.4355
- Hanke F, Mu ¨ller B, Wongwathanarat A, Marek A, Janka HT (2013) SASI activity in three-dimensional neutrino-hydrodynamics simulations of supernova cores. Astrophys J 770:66. https://doi.org/10. 1088/0004-637X/770/1/66. arXiv:1303.6269
- Hawley J, Blondin J, Lindahl G, Lufkin E (2012) VH-1: multidimensional ideal compressible hydrodynamics code. Astrophysics Source Code Library ascl:1204.007. http://www.ascl.net/1204. 007
- Hecht J, Ofer D, Alon U, Shvarts D, Orszag SA, Shvarts D, McCrory RL (1995) Three-dimensional simulations and analysis of the nonlinear stage of the Rayleigh-Taylor instability. Laser Part Beams 13:423. https://doi.org/10.1017/S026303460000954X
- Heger A, Woosley SE (2010) Nucleosynthesis and evolution of massive metal-free stars. Astrophys J 724:341-373. https://doi.org/10.1088/0004-637X/724/1/341
- Heger A, Woosley SE, Spruit HC (2005) Presupernova evolution of differentially rotating massive stars including magnetic fields. Astrophys J 626:350-363. https://doi.org/10.1086/429868. arXiv:astroph/0409422
- Herant M, Benz W (1991) Hydrodynamical instabilities and mixing in SN 1987A: two-dimensional simulations of the first 3 months. Astrophys J 370:L81. https://doi.org/10.1086/185982
- Herant M, Benz W (1992) Postexplosion hydrodynamics of SN 1987A. Astrophys J 387:294. https://doi. org/10.1086/171081
- Herant M, Benz W, Colgate S (1992) Postcollapse hydrodynamics of SN 1987A: two-dimensional simulations of the early evolution. Astrophys J 395:642-653. https://doi.org/10.1086/171685
- Herant M, Benz W, Hix WR, Fryer CL, Colgate SA (1994) Inside the supernova: a powerful convective engine. Astrophys J 435:339-361. https://doi.org/10.1086/174817. arXiv:astro-ph/9404024
- Hix WR, Meyer BS (2006) Thermonuclear kinetics in astrophysics. Nucl Phys A 777:188-207. https:// doi.org/10.1016/j.nuclphysa.2004.10.009
- Hix WR, Parete-Koon ST, Freiburghaus C, Thielemann FK (2007) The QSE-reduced nuclear reaction network for silicon burning. Astrophys J 667(1):476-488. https://doi.org/10.1086/520672
- Hobbs G, Lorimer DR, Lyne AG, Kramer M (2005) A statistical study of 233 pulsar proper motions. Mon Not R Astron Soc 360:974-992. https://doi.org/10.1111/j.1365-2966.2005.09087.x
- Hockney RW (1965) A fast direct solution of Poisson's equation using fourier analysis. J ACM 12:95
- Hotta H, Rempel M, Yokoyama T (2015) Efficient small-scale dynamo in the solar convection zone.
- Astrophys J 803(1):42. https://doi.org/10.1088/0004-637X/803/1/42
- Houck JC, Chevalier RA (1992) Linear stability analysis of spherical accretion flows onto compact objects. Astrophys J 395:592. https://doi.org/10.1086/171679

<!-- image -->

- Huang K, Wu H, Yu H, Yan D (2011) Cures for numerical shock instability in HLLC solver. Int J Numer Meth Fluids 65(9):1026-1038. https://doi.org/10.1002/fld.2217
- Hu ¨depohl L, Mu ¨ller B, Janka HT, Marek A, Raffelt GG (2010) Neutrino signal of electron-capture supernovae from core collapse to cooling. Phys Rev Lett 104(25):251101. https://doi.org/10.1103/ PhysRevLett.104.251101. arXiv:0912.0260
- Huete C, Abdikamalov E (2019) Response of nuclear-dissociating shocks to vorticity perturbations. Phys Scripta 94(9):094002. https://doi.org/10.1088/1402-4896/ab0228
- Huete C, Abdikamalov E, Radice D (2018) The impact of vorticity waves on the shock dynamics in corecollapse supernovae. Mon Not R Astron Soc 475(3):3305-3323. https://doi.org/10.1093/mnras/ stx3360
- Hungerford AL, Fryer CL, Warren MS (2003) Gamma-ray lines from asymmetric supernovae. Astrophys J 594(1):390-403. https://doi.org/10.1086/376776
- Iliadis C (2007) Nuclear physics of stars. Wiley, Weinheim
- Imshennik VS, Nadezhin DK (1972) Neutrino thermal conductivity in collapsing stars. Zh Eksp Teor Fiz 63:1548-1561
- Isenberg JA (1978) Waveless approximation theories of gravities. University of Maryland Preprint. arXiv: gr-qc/0702113
- Iwakami W, Kotake K, Ohnishi N, Yamada S, Sawada K (2008) Three-dimensional simulations of standing accretion shock instability in core-collapse supernovae. Astrophys J 678:1207-1222. https://doi.org/10.1086/533582. arXiv:0710.2191
- Iwakami W, Kotake K, Ohnishi N, Yamada S, Sawada K (2009) Effects of rotation on standing accretion shock instability in nonlinear phase for core-collapse supernovae. Astrophys J 700:232-242. https:// doi.org/10.1088/0004-637X/700/1/232. arXiv:0811.0651
- Iwakami W, Nagakura H, Yamada S (2014) Critical surface for explosions of rotational core-collapse supernovae. Astrophys J 793(1):5. https://doi.org/10.1088/0004-637X/793/1/5
- Janka HT (2001) Conditions for shock revival by neutrino heating in core-collapse supernovae. Astron Astrophys 368:527-560. https://doi.org/10.1051/0004-6361:20010012. arXiv:astro-ph/0008432
- Janka HT (2012) Explosion mechanisms of core-collapse supernovae. Annu Rev Nucl Part Sci 62:407-451. https://doi.org/10.1146/annurev-nucl-102711-094901. arXiv:1206.2503
- Janka HT (2017) Neutron star kicks by the gravitational tug-boat mechanism in asymmetric supernova explosions: progenitor and explosion dependence. Astrophys J 837:84. https://doi.org/10.3847/15384357/aa618e
- Janka HT, Mu ¨ller E (1994) Neutron star recoils from anisotropic supernovae. Astron Astrophys 290:496-502
- Janka HT, Mu ¨ller E (1995) The first second of a type II supernova: convection, accretion, and shock propagation. Astrophys J 448:L109-L113. https://doi.org/10.1086/309604
- Janka HT, Mu ¨ller E (1996) Neutrino heating, convection, and the mechanism of type-II supernova explosions. Astron Astrophys 306:167-198
- Janka HT, Zwerger T, Mo ¨nchmeyer R (1993) Does artificial viscosity destroy prompt type-II supernova explosion. Astron Astrophys 268:360-368
- Janka HT, Hanke F, Hu ¨depohl L, Marek A, Mu ¨ller B, Obergaulinger MM (2012) Core-collapse supernovae: reflections and directions. Prog Theor Exp Phys 1:01A309. https://doi.org/10.1093/ ptep/pts067. arXiv:1211.1378
- Janka HT, Melson T, Summa A (2016) Physics of core-collapse supernovae in three dimensions: a sneak preview. Annu Rev Nucl Part Sci 66:341-375. https://doi.org/10.1146/annurev-nucl-102115-044747
- Jerkstrand A, Timmes FX, Magkotsios G, Sim SA, Fransson C, Spyromilio J, Mu ¨ller B, Heger A, Sollerman J, Smartt SJ (2015) Constraints on explosive silicon burning in core-collapse supernovae from measured Ni/Fe ratios. Astrophys J 807:110. https://doi.org/10.1088/0004-637X/807/1/110
- Joggerst CC, Woosley SE, Heger A (2009) Mixing in zero- and solar-metallicity supernovae. Astrophys J 693(2):1780-1802. https://doi.org/10.1088/0004-637X/693/2/1780
- Joggerst CC, Almgren A, Bell J, Ae H, Whalen D, Woosley SE (2010a) The nucleosynthetic imprint of 15-40 M /C12 primordial supernovae on metal-poor stars. Astrophys J 709(1):11-26. https://doi.org/10. 1088/0004-637X/709/1/11
- Joggerst CC, Almgren A, Woosley SE (2010b) Three-dimensional simulations of Rayleigh-Taylor mixing in core-collapse supernovae with castro. Astrophys J 723(1):353-363. https://doi.org/10. 1088/0004-637X/723/1/353

<!-- image -->

- Johnston S, Hobbs G, Vigeland S, Kramer M, Weisberg JM, Lyne AG (2005) Evidence for alignment of the rotation and velocity vectors in pulsars. Mon Not R Astron Soc 364:1397-1412. https://doi.org/ 10.1111/j.1365-2966.2005.09669.x
- Jones S, Hirschi R, Nomoto K, Fischer T, Timmes FX, Herwig F, Paxton B, Toki H, Suzuki T, Martı ´nezPinedo G, Lam YH, Bertolli MG (2013) Advanced burning stages and fate of 8-10 M /C12 stars. Astrophys J 772:150. https://doi.org/10.1088/0004-637X/772/2/150
- Jones S, Ro ¨pke FK, Pakmor R, Seitenzahl IR, Ohlmann ST, Edelmann PVF (2016) Do electron-capture supernovae make neutron stars? First multidimensional hydrodynamic simulations of the oxygen deflagration. Astron Astrophys 593:A72. https://doi.org/10.1051/0004-6361/201628321. arXiv: 1602.05771
- Jones S, Andrassy R, Sandalski S, Davis A, Woodward P, Herwig F (2017) Idealized hydrodynamic simulations of turbulent oxygen-burning shell convection in 4 p geometry. Mon Not R Astron Soc 465(3):2991-3010. https://doi.org/10.1093/mnras/stw2783
- Just O, Obergaulinger M, Janka HT (2015) A new multidimensional, energy-dependent two-moment transport code for neutrino-hydrodynamics. Mon Not R Astron Soc 453:3386-3413. https://doi.org/ 10.1093/mnras/stv1892
- Just O, Bollig R, Janka HT, Obergaulinger M, Glas R, Nagataki S (2018) Core-collapse supernova simulations in one and two dimensions: comparison of codes and approximations. Mon Not R Astron Soc 481(4):4786-4814. https://doi.org/10.1093/mnras/sty2578
- Kageyama A, Sato T (2004) ''Yin-Yang grid'': an overset grid in spherical geometry. Geochem Geophys Geosyst 5(9):Q09005. https://doi.org/10.1029/2004GC000734
- Kalogera V, Bizouard MA, Burrows A, Janka HT, Kotake K, Messer B, Mezzacappa T, Mueller B, Mueller E, Papa MA, Reddy S, Rosswog S (2019) The yet-unobserved multi-messenger gravitational-wave universe. Bull Am Astron Soc 51(3):239
- Kane J, Arnett D, Remington BA, Glendinning SG, Baza ´n G, Mu ¨ller E, Fryxell BA, Teyssier R (2000) Two-dimensional versus three-dimensional supernova hydrodynamic instability growth. Astrophys J 528:989-994. https://doi.org/10.1086/308220
- Ka ¨ppeli R, Mishra S (2016) A well-balanced finite volume scheme for the Euler equations with gravitation. The exact preservation of hydrostatic equilibrium with arbitrary entropy stratification. Astron Astrophys 587:A94. https://doi.org/10.1051/0004-6361/201527815
- Ka ¨ppeli R, Whitehouse SC, Scheidegger S, Pen UL, Liebendo ¨rfer M (2011) FISH: a three-dimensional parallel magnetohydrodynamics code for astrophysical applications. Astrophys J Suppl 195(2):20. https://doi.org/10.1088/0067-0049/195/2/20
- Kasen D, Woosley SE (2009) Type II supernovae: model light curves and standard candle relationships. Astrophys J 703:2205-2216. https://doi.org/10.1088/0004-637X/703/2/2205
- Kastaun W (2006) High-resolution shock capturing scheme for ideal hydrodynamics in general relativity optimized for quasistationary solutions. Phys Rev D 74(12):124024. https://doi.org/10.1103/ PhysRevD.74.124024
- Kazeroni R, Guilet J, Foglizzo T (2016) New insights on the spin-up of a neutron star during core collapse. Mon Not R Astron Soc 456:126-135. https://doi.org/10.1093/mnras/stv2666
- Kazeroni R, Guilet J, Foglizzo T (2017) Are pulsars spun up or down by SASI spiral modes? Mon Not R Astron Soc 471(1):914-925. https://doi.org/10.1093/mnras/stx1566
- Kazeroni R, Krueger BK, Guilet J, Foglizzo T, Pomare `de D (2018) The non-linear onset of neutrinodriven convection in two- and three-dimensional core-collapse supernovae. Mon Not R Astron Soc 480(1):261-280. https://doi.org/10.1093/mnras/sty1742
- Kazeroni R, Abdikamalov E (2019) The impact of progenitor asymmetries on the neutrino-driven convection in core-collapse supernovae. arXiv:1911.08819
- Keil W, Janka HT, Mu ¨ller E (1996) Ledoux convection in protoneutron stars--a clue to supernova nucleosynthesis? Astrophys J 473:L111-L114. https://doi.org/10.1086/310404. arXiv:astro-ph/ 9610203
- Keller SC, Bessell MS, Frebel A, Casey AR, Asplund M, Jacobson HR, Lind K, Norris JE, Yong D, Heger A, Magic Z, da Costa GS, Schmidt BP, Tisserand P (2014) A single low-energy, iron-poor supernova as the source of metals in the star SMSS J031300.36-670839.3. Nature 506(7489):463-466. https://doi.org/10.1038/nature12990
- Kifonidis K, Plewa T, Janka HT, Mu ¨ller E (2000) Nucleosynthesis and clump formation in a corecollapse supernova. Astrophys J 531:L123-L126. https://doi.org/10.1086/312541. arXiv:astro-ph/ 9911183

<!-- image -->

- Kifonidis K, Plewa T, Janka HT, Mu ¨ller E (2003) Non-spherical core collapse supernovae. I. Neutrinodriven convection, Rayleigh-Taylor instabilities, and the formation and propagation of metal clumps. Astron Astrophys 408:621-649. https://doi.org/10.1051/0004-6361:20030863. arXiv:astroph/0302239
- Kifonidis K, Plewa T, Scheck L, Janka HT, Mu ¨ller E (2006) Non-spherical core collapse supernovae. II. The late-time evolution of globally anisotropic neutrino-driven explosions and their implications for SN 1987 A. Astron Astrophys 453:661-678. https://doi.org/10.1051/0004-6361:20054512. arXiv: astro-ph/0511369
- Kippenhahn R, Weigert A, Weiss A (2012) Stellar structure and evolution, 2nd edn. Springer, Berlin. https://doi.org/10.1007/978-3-642-30304-3
- Kirsebom OS, Jones S, Stro ¨mberg DF, Martı ´nez-Pinedo G, Langanke K, Roepke FK, Brown BA, Eronen T, Fynbo HOU, Hukkanen M, Idini A, Jokinen A, Kankainen A, Kostensalo J, Moore I, Mo ¨ller H, Ohlmann ST, Penttila ¨ H, Riisager K, Rinta-Antila S, Srivastava PC, Suhonen J, Trzaska WH, A ¨ ysto ¨ J (2019) Discovery of an exceptionally strong b -decay transition of 20 F and implications for the fate of intermediate-mass stars. Phys Rev Lett 123(26):262701. https://doi.org/10.1103/PhysRevLett. 123.262701. arXiv:1905.09407
- Kitaura FS, Janka HT, Hillebrandt W (2006) Explosions of O-Ne-Mg cores, the crab supernova, and subluminous type II-P supernovae. Astron Astrophys 450:345-350. https://doi.org/10.1051/00046361:20054703. arXiv:astro-ph/0512065
- Koldoba AV, Romanova MM, Ustyugova GV, Lovelace RVE (2002) Three-dimensional magnetohydrodynamic simulations of accretion to an inclined rotator: the ''cubed sphere'' method. Astrophys J 576(1):L53-L56. https://doi.org/10.1086/342780
- Kotake K (2013) Multiple physical elements to determine the gravitational-wave signatures of corecollapse supernovae. C R Physique 14:318-351. https://doi.org/10.1016/j.crhy.2013.01.008
- Kotake K, Yamada S, Sato K (2003) Anisotropic neutrino radiation in rotational core collapse. Astrophys J 595(1):304-316. https://doi.org/10.1086/377196
- Kotake K, Sato K, Takahashi K (2006) Explosion mechanism, neutrino burst and gravitational wave in core-collapse supernovae. Rep Progr Phys 69(4):971-1143. https://doi.org/10.1088/0034-4885/69/4/ R03. arXiv:astro-ph/0509456
- Kovalenko IG, Eremin MA (1998) Instability of spherical accretion: I. Shock-free Bondi accretion. Mon Not R Astron Soc 298:861-870. https://doi.org/10.1046/j.1365-8711.1998.01667.x
- Kozma C, Fransson C (1998) Late spectral evolution of SN 1987A. II. Line emission. Astrophys J 497(1):431-457. https://doi.org/10.1086/305452
- Kraichnan RH (1967) Inertial ranges in two-dimensional turbulence. Phys Fluids 10:1417-1423. https:// doi.org/10.1063/1.1762301
- Kuhlen M, Woosley WE, Glatzmaier GA (2003) 3D anelastic simulations of convection in massive stars. In: Turcotte S, Keller SC, Cavallo RM (eds) 3D stellar evolution, ASP conference series, vol 293. Astronomical Society of the Pacific, San Francisco, p 147. arXiv:astro-ph/0210557
- Kuroda T, Kotake K, Takiwaki T (2012) Fully general relativistic simulations of core-collapse supernovae with an approximate neutrino transport. Astrophys J 755:11. https://doi.org/10.1088/ 0004-637X/755/1/11. arXiv:1202.2487
- Kuroda T, Kotake K, Takiwaki T (2016a) A new gravitational-wave signature from standing accretion shock instability in supernovae. Astrophys J 829(1):L14. https://doi.org/10.3847/2041-8205/829/1/ L14
- Kuroda T, Takiwaki T, Kotake K (2016b) A New multi-energy neutrino radiation-hydrodynamics code in full general relativity and its application to the gravitational collapse of massive stars. Astrophys J Suppl 222:20. https://doi.org/10.3847/0067-0049/222/2/20
- Kuroda T, Kotake K, Takiwaki T, Thielemann FK (2018) A full general relativistic neutrino radiationhydrodynamics simulation of a collapsing very massive star and the formation of a black hole. Mon Not R Astron Soc 477:L80-L84. https://doi.org/10.1093/mnrasl/sly059
- Lai D, Goldreich P (2000) Growth of perturbations in gravitational collapse and accretion. Astrophys J 535:402-411. https://doi.org/10.1086/308821. arXiv:astro-ph/9906400
- Lai D, Chernoff DF, Cordes JM (2001) Pulsar jets: implications for neutron star kicks and initial spins. Astrophys J 549:1111-1118. https://doi.org/10.1086/319455
- Laming JM (2007) Analytic approach to the stability of standing accretion shocks: application to corecollapse supernovae. Astrophys J 659:1449-1457. https://doi.org/10.1086/512534. arXiv:astro-ph/ 0701264

<!-- image -->

- Laming JM (2008) Erratum: ''Analytic approach to the stability of standing accretion shocks: application to core-collapse supernovae'' (ApJ, 659, 1449 [2007]). Astrophys J 687(2):1461-1463. https://doi. org/10.1086/592088
- Laney CB (1998) Computational gasdynamics. Cambridge University Press, Cambridge. https://doi.org/ 10.1017/CBO9780511605604
- Lattimer JM (2012) The nuclear equation of state and neutron star masses. Annu Rev Nucl Part Sci 62(1):485-515. https://doi.org/10.1146/annurev-nucl-102711-095018
- Lecoanet D, Schwab J, Quataert E, Bildsten L, Timmes FX, Burns KJ, Vasil GM, Oishi JS, Brown BP (2016) Turbulent chemical diffusion in convectively bounded carbon flames. Astrophys J 832(1):71. https://doi.org/10.3847/0004-637X/832/1/71
- Lehner L, Pretorius F (2014) Numerical relativity and astrophysics. Annu Rev Astron Astrophys 52:661-694. https://doi.org/10.1146/annurev-astro-081913-040031
- Lentz EJ, Mezzacappa A, Bronson Messer OE, Liebendo ¨rfer M, Hix WR, Bruenn SW (2012) On the requirements for realistic modeling of neutrino transport in simulations of core-collapse supernovae. Astrophys J 747:73. https://doi.org/10.1088/0004-637X/747/1/73. arXiv:1112.3595
- Lentz EJ, Bruenn SW, Hix WR, Mezzacappa A, Messer OEB, Endeve E, Blondin JM, Harris JA, Marronetti P, Yakunin KN (2015) Three-dimensional core-collapse supernova simulated using a 15 M /C12 progenitor. Astrophys J 807:L31. https://doi.org/10.1088/2041-8205/807/2/L31
- Leung SC, Nomoto K (2019) Final evolution of super-AGB stars and supernovae triggered by electron capture. Publ Astron Soc Australia 36:e006. https://doi.org/10.1017/pasa.2018.49
- Leung SC, Nomoto K, Suzuki T (2020) Electron-capture supernovae of super-AGB stars: sensitivity on input physics. Astrophys J 889(1):34. https://doi.org/10.3847/1538-4357/ab5d2f. arXiv:1901.11438
- LeVeque RJ (1998a) Balancing source terms and flux gradients in high-resolution Godunov methods: the quasi-steady wave-propagation algorithm. J Comput Phys 146(1):346-365. https://doi.org/10.1006/ jcph.1998.6058
- LeVeque RJ (1998b) Nonlinear conservation laws and finite volume methods. In: Steiner O, Gautschy A (eds) Computational methods for astrophysical fluid flow. Saas-Fee Advanced Course, vol 27. Springer, Berlin, pp 1-159. https://doi.org/10.1007/3-540-31632-9\_1
- Liebendo ¨rfer M, Mezzacappa A, Thielemann FK, Messer OE, Hix WR, Bruenn SW (2001) Probing the gravitational well: no supernova explosion in spherical symmetry with general relativistic Boltzmann neutrino transport. Phys Rev D 63(10):103004:1-13. https://doi.org/10.1103/ PhysRevD.63.103004. arXiv:astro-ph/0006418
- Liebendo ¨rfer M, Messer OEB, Mezzacappa A, Bruenn SW, Cardall CY, Thielemann FK (2004) A finite difference representation of neutrino radiation hydrodynamics in spherically symmetric general relativistic spacetime. Astrophys J Suppl 150:263-316. https://doi.org/10.1086/380191. arXiv:astroph/0207036
- Liebendo ¨rfer M, Rampp M, Janka HT, Mezzacappa A (2005) Supernova simulations with Boltzmann neutrino transport: a comparison of methods. Astrophys J 620:840-860. https://doi.org/10.1086/ 427203
- Liou MS (2000) Mass flux schemes and connection to shock instability. J Comput Phys 160(2):623-648. https://doi.org/10.1006/jcph.2000.6478
- Livne E (1993) An implicit method for two-dimensional hydrodynamics. Astrophys J 412:634. https:// doi.org/10.1086/172950
- Livne E, Burrows A, Walder R, Lichtenstadt I, Thompson TA (2004) Two-dimensional, time-dependent, multigroup, multiangle radiation hydrodynamics test simulation in the core-collapse supernova context. Astrophys J 609:277-287. https://doi.org/10.1086/421012. arXiv:astro-ph/0312633
- Lucy LB (1977) A numerical approach to the testing of the fission hypothesis. Astron J 82:1013-1024. https://doi.org/10.1086/112164
- Lund T, Marek A, Lunardini C, Janka HT, Raffelt G (2010) Fast time variations of supernova neutrino fluxes and their detectability. Phys Rev D 82(6):063007. https://doi.org/10.1103/PhysRevD.82. 063007. arXiv:1006.1889
- Mabanta QA, Murphy JW (2018) How turbulence enables core-collapse supernova explosions. Astrophys J 856:22. https://doi.org/10.3847/1538-4357/aaaec7
- Maeda K, Nomoto K (2003) Bipolar supernova explosions: nucleosynthesis and implications for abundances in extremely metal-poor stars. Astrophys J 598(2):1163-1200. https://doi.org/10.1086/ 378948

<!-- image -->

- Mao J, Ono M, Nagataki S, Ma H, Ito H, Matsumoto J, Dainotti MG, Lee SH (2015) Matter mixing in core-collapse supernova ejecta: large density perturbations in the progenitor star? Astrophys J 808(2):164. https://doi.org/10.1088/0004-637X/808/2/164
- Marek A, Janka HT (2009) Delayed neutrino-driven supernova explosions aided by the standing accretion-shock instability. Astrophys J 694:664-696. https://doi.org/10.1088/0004-637X/694/1/ 664. arXiv:0708.3372
- Marek A, Janka HT, Buras R, Liebendo ¨rfer M, Rampp M (2005) On ion-ion correlation effects during stellar core collapse. Astron Astrophys 443:201-210. https://doi.org/10.1051/0004-6361:20053236. arXiv:astro-ph/0504291
- Marek A, Janka HT, Mu ¨ller E (2009) Equation-of-state dependent features in shock-oscillation modulated neutrino and gravitational-wave signals from supernovae. Astron Astrophys 496:475-494. https:// doi.org/10.1051/0004-6361/200810883. arXiv:0808.4136
- Martı ´ JM, Mu ¨ller E (2015) Grid-based methods in relativistic hydrodynamics and magnetohydrodynamics. Living Rev Comput Astrophys 1:3. https://doi.org/10.1007/lrca-2015-3
- Matzner CD, McKee CF (1999) The expulsion of stellar envelopes in core-collapse supernovae. Astrophys J 510:379-403. https://doi.org/10.1086/306571. arXiv:astro-ph/9807046
- Mazurek TJ (1982) The energetics of adiabatic shocks in stellar collapse. Astrophys J 259:L13-L17. https://doi.org/10.1086/183839
- McCray R (1993) Supernova 1987A revisited. Annu Rev Astron Astrophys 31:175-216. https://doi.org/ 10.1146/annurev.aa.31.090193.001135
- Meakin CA, Arnett D (2006) Active carbon and oxygen shell burning hydrodynamics. Astrophys J 637:L53-L56. https://doi.org/10.1086/500544. arXiv:astro-ph/0601348
- Meakin CA, Arnett D (2007a) Anelastic and compressible simulations of stellar oxygen burning. Astrophys J 665:690-697. https://doi.org/10.1086/519372. arXiv:astro-ph/0611317
- Meakin CA, Arnett D (2007b) Turbulent convection in stellar interiors. I. Hydrodynamic simulation. Astrophys J 667:448-475. https://doi.org/10.1086/520318. arXiv:astro-ph/0611315
- Meakin CA, Arnett WD (2010) Some properties of the kinetic energy flux and dissipation in turbulent stellar convection zones. Astrophys Space Sci 328:221-225. https://doi.org/10.1007/s10509-0100301-6
- Melson T (2013) Core-collapse supernova hydrodynamics on the Yin-Yang grid with PROMETHEUSVERTEX. Master's thesis, Ludwig-Maximilians Universtia ¨t Mu ¨nchen
- Melson T, Janka HT, Bollig R, Hanke F, Marek A, Mu ¨ller B (2015a) Neutrino-driven explosion of a 20 solar-mass star in three dimensions enabled by strange-quark contributions to neutrino-nucleon scattering. Astrophys J 808:L42. https://doi.org/10.1088/2041-8205/808/2/L42
- Melson T, Janka HT, Marek A (2015b) Neutrino-driven supernova of a low-mass iron-core progenitor boosted by three-dimensional turbulent convection. Astrophys J 801:L24. https://doi.org/10.1088/ 2041-8205/801/2/L24
- Melson T, Kresse D, Janka HT (2020) Resolution study for three-dimensional supernova simulations with the PROMETHEUS-VERTEX code. Astrophys J 891(1):27. https://doi.org/10.3847/1538-4357/ ab72a7. arXiv:1904.01699
- Mezzacappa A (2005) Ascertaining the core collapse supernova mechanism: the state of the art and the road ahead. Annu Rev Nucl Part Sci 55:467-515. https://doi.org/10.1146/annurev.nucl.55.090704. 151608
- Mezzacappa A (2020) In prep., Living Rev Comput Astrophys
- Mezzacappa A, Calder AC, Bruenn SW, Blondin JM, Guidry MW, Strayer MR, Umar AS (1998) An investigation of neutrino-driven convection and the core collapse supernova mechanism using multigroup neutrino transport. Astrophys J 495:911-926. https://doi.org/10.1086/305338. arXiv: astro-ph/9709188
- Michel A (2019) Modeling of silicon burning during late stages of stellar evolution. PhD thesis, Universita ¨t Heidelberg
- Miczek F, Ro ¨pke FK, Edelmann PVF (2015) New numerical solver for flows at various Mach numbers. Astron Astrophys 576:A50. https://doi.org/10.1051/0004-6361/201425059
- Mignone A, Bodo G (2005) An HLLC Riemann solver for relativistic flows : I. Hydrodynamics. Mon Not R Astron Soc 364:126-136. https://doi.org/10.1111/j.1365-2966.2005.09546.x
- Mirizzi A, Tamborra I, Janka HT, Saviano N, Scholberg K, Bollig R, Hu ¨depohl L, Chakraborty S (2016) Supernova neutrinos: production, oscillations and detection. Riv Nuovo Cimento 39:1-112. https:// doi.org/10.1393/ncr/i2016-10120-8

<!-- image -->

- Moca ´k M, Mu ¨ller E, Weiss A, Kifonidis K (2009) The core helium flash revisited. II. Two and threedimensional hydrodynamic simulations. Astron Astrophys 501:659-677. https://doi.org/10.1051/ 0004-6361/200811414
- Moca ´k M, Meakin C, Viallet M, Arnett D (2014) Compressible hydrodynamic mean-field equations in spherical geometry and their application to turbulent stellar convection data. arXiv:1401.5176
- Moca ´k M, Meakin C, Campbell SW, Arnett WD (2018) Turbulent mixing and nuclear burning in stellar interiors. Mon Not R Astron Soc 481(3):2918-2932. https://doi.org/10.1093/mnras/sty2392
- Mocz P, Vogelsberger M, Sijacki D, Pakmor R, Hernquist L (2014) A discontinuous Galerkin method for solving the fluid and magnetohydrodynamic equations in astrophysical simulations. Mon Not R Astron Soc 437(1):397-414. https://doi.org/10.1093/mnras/stt1890
- Morozova V, Radice D, Burrows A, Vartanyan D (2018) The gravitational wave signal from corecollapse supernovae. Astrophys J 861(1):10. https://doi.org/10.3847/1538-4357/aac5f1
- Mo ¨sta P, Mundim BC, Faber JA, Haas R, Noble SC, Bode T, Lo ¨ffler F, Ott CD, Reisswig C, Schnetter E (2014a) GRHydro: a new open-source general-relativistic magnetohydrodynamics code for the Einstein toolkit. Class Quantum Grav 31(1):015005. https://doi.org/10.1088/0264-9381/31/1/ 015005
- Mo ¨sta P, Richers S, Ott CD, Haas R, Piro AL, Boydstun K, Abdikamalov E, Reisswig C, Schnetter E (2014b) Magnetorotational core-collapse supernovae in three dimensions. Astrophys J 785:L29. https://doi.org/10.1088/2041-8205/785/2/L29
- Mo ¨sta P, Roberts LF, Halevi G, Ott CD, Lippuner J, Haas R, Schnetter E (2018) r-process nucleosynthesis from three-dimensional magnetorotational core-collapse supernovae. Astrophys J 864(2):171. https://doi.org/10.3847/1538-4357/aad6ec
- Mu ¨ller E (1994) Fundamentals of gas-dynamical simulations. In: Contopoulos G, Spyrou NK, Vlahos L (eds) Galactic dynamics and N-body simulations. Lecture Notes in Physics, vol 433, Springer, Berlin, pp 313-363. https://doi.org/10.1007/3-540-57983-4\_23
- Mu ¨ller E (1998) Simulation of astrophysical fluid flow. In: Steiner O, Gautschy A (eds) Computational methods for astrophysical fluid flow. Saas-Fee Advanced Course, vol 27. Springer, Berlin, pp 343-494. https://doi.org/10.1007/3-540-31632-9\_4
- Mu ¨ller B (2009) Multi-dimensional relativistic simulations of core-collapse supernovae with energydependent neutrino transport. PhD thesis, Technische Universita ¨t Mu ¨nchen. http://mediatum.ub.tum. de/?id=800389
- Mu ¨ller B (2015) The dynamics of neutrino-driven supernova explosions after shock revival in 2D and 3D. Mon Not R Astron Soc 453:287-310. https://doi.org/10.1093/mnras/stv1611
- Mu ¨ller B (2016) The status of multi-dimensional core-collapse supernova models. Publ Astron Soc Australia 33:e048. https://doi.org/10.1017/pasa.2016.40
- Mu ¨ller B (2019a) A critical assessment of turbulence models for 1D core-collapse supernova simulations. Mon Not R Astron Soc 487(4):5304-5323. https://doi.org/10.1093/mnras/stz1594
- Mu ¨ller B (2019b) Neutrino emission as diagnostics of core-collapse supernovae. Annu Rev Nucl Part Sci 69(1):annurev. https://doi.org/10.1146/annurev-nucl-101918-023434
- Mu ¨ller B, Chan C (2019) An FFT-based solution method for the poisson equation on 3D spherical polar grids. Astrophys J 870(1):43. https://doi.org/10.3847/1538-4357/aaf100
- Mu ¨ller B, Janka HT (2014) A new multi-dimensional general relativistic neutrino hydrodynamics code for core-collapse supernovae. IV. The neutrino signal. Astrophys J 788:82. https://doi.org/10.1088/ 0004-637X/788/1/82
- Mu ¨ller B, Janka HT (2015) Non-radial instabilities and progenitor asphericities in core-collapse supernovae. Mon Not R Astron Soc 448:2141-2174. https://doi.org/10.1093/mnras/stv101
- Mu ¨ller E, Steinmetz M (1995) Simulating self-gravitating hydrodynamic flows. Comput Phys Commun 89:45-58. https://doi.org/10.1016/0010-4655(94)00185-5. arXiv:astro-ph/9402070
- Mu ¨ller E, Fryxell B, Arnett D (1991) Instability and clumping in SN 1987A. Astron Astrophys 251:505-514
- Mu ¨ller B, Dimmelmeier H, Mu ¨ller E (2008) Exploring the relativistic regime with Newtonian hydrodynamics. II. An effective gravitational potential for rapid rotation. Astron Astrophys 489:301-314. https://doi.org/10.1051/0004-6361:200809609. arXiv:0802.2459
- Mu ¨ller B, Janka HT, Dimmelmeier H (2010) A new multi-dimensional general relativistic neutrino hydrodynamic code for core-collapse supernovae. I. Method and code tests in spherical symmetry. Astrophys J Suppl 189:104-133. https://doi.org/10.1088/0067-0049/189/1/104. arXiv:1001.4841

<!-- image -->

- Mu ¨ller B, Janka HT, Heger A (2012a) New two-dimensional models of supernova explosions by the neutrino-heating mechanism: evidence for different instability regimes in collapsing stellar cores. Astrophys J 761:72. https://doi.org/10.1088/0004-637X/761/1/72
- Mu ¨ller B, Janka HT, Marek A (2012b) A new multi-dimensional general relativistic neutrino hydrodynamics code for core-collapse supernovae. II. Relativistic explosion models of corecollapse supernovae. Astrophys J 756:84. https://doi.org/10.1088/0004-637X/756/1/84
- Mu ¨ller B, Janka HT, Marek A (2013) A new multi-dimensional general relativistic neutrino hydrodynamics code of core-collapse supernovae. III. Gravitational wave signals from supernova explosion models. Astrophys J 766:43. https://doi.org/10.1088/0004-637X/766/1/43. arXiv:1210. 6984
- Mu ¨ller B, Heger A, Liptai D, Cameron JB (2016a) A simple approach to the supernova progenitorexplosion connection. Mon Not R Astron Soc 460:742-764. https://doi.org/10.1093/mnras/stw1083
- Mu ¨ller B, Viallet M, Heger A, Janka HT (2016b) The last minutes of oxygen shell burning in a massive star. Astrophys J 833:124. https://doi.org/10.3847/1538-4357/833/1/124
- Mu ¨ller B, Melson T, Heger A, Janka HT (2017a) Supernova simulations from a 3D progenitor model: impact of perturbations and evolution of explosion properties. Mon Not R Astron Soc 472:491-513. https://doi.org/10.1093/mnras/stx1962
- Mu ¨ller T, Prieto JL, Pejcha O, Clocchiatti A (2017b) The nickel mass distribution of normal type II supernovae. Astrophys J 841:127. https://doi.org/10.3847/1538-4357/aa72f1
- Mu ¨ller B, Gay DW, Heger A, Tauris TM, Sim SA (2018) Multidimensional simulations of ultrastripped supernovae to shock breakout. Mon Not R Astron Soc 479:3675-3689. https://doi.org/10.1093/ mnras/sty1683
- Mu ¨ller B, Tauris TM, Heger A, Banerjee P, Qian YZ, Powell J, Chan C, Gay DW, Langer N (2019) Three-dimensional simulations of neutrino-driven core-collapse supernovae from low-mass single and binary star progenitors. Mon Not R Astron Soc 484:3307-3324. https://doi.org/10.1093/mnras/ stz216
- Murphy JW, Burrows A (2008) Criteria for core-collapse supernova explosions by the neutrino mechanism. Astrophys J 688:1159-1175. https://doi.org/10.1086/592214. arXiv:0805.3345
- Murphy JW, Meakin C (2011) A global turbulence model for neutrino-driven convection in core-collapse supernovae. Astrophys J 742:74. https://doi.org/10.1088/0004-637X/742/2/74. arXiv:1106.5496
- Murphy JW, Dolence JC, Burrows A (2013) The dominance of neutrino-driven convection in corecollapse supernovae. Astrophys J 771:52. https://doi.org/10.1088/0004-637X/771/1/52
- Murphy JW, Mabanta Q, Dolence JC (2019) A comparison of explosion energies for simulated and observed core-collapse supernovae. Mon Not R Astron Soc 489(1):641-652. https://doi.org/10. 1093/mnras/stz2123
- Nagakura H, Sumiyoshi K, Yamada S (2019) Possible early linear acceleration of proto-neutron stars via asymmetric neutrino emission in core-collapse supernovae. Astrophys J 880(2):L28. https://doi.org/ 10.3847/2041-8213/ab30ca
- Nagataki S (2000) Effects of jetlike explosion in SN 1987A. Astrophys J Suppl 127(1):141-157. https:// doi.org/10.1086/313317
- Nagataki S, Shimizu TM, Sato K (1998) Matter mixing from axisymmetric supernova explosion. Astrophys J 495(1):413-423. https://doi.org/10.1086/305258
- Nakamura K, Kuroda T, Takiwaki T, Kotake K (2014) Impacts of rotation on three-dimensional hydrodynamics of core-collapse supernovae. Astrophys J 793:45. https://doi.org/10.1088/0004637X/793/1/45
- Nakamura K, Takiwaki T, Kuroda T, Kotake K (2015) Systematic features of axisymmetric neutrinodriven core-collapse supernova models in multiple progenitors. Publ Astron Soc Japan 67:107. https://doi.org/10.1093/pasj/psv073
- Nakamura K, Takiwaki T, Kotake K (2019) Long-term simulations of multi-dimensional core-collapse supernovae: implications for neutron star kicks. Publ Astron Soc Japan 71(5):98. https://doi.org/10. 1093/pasj/psz080
- Ng CY, Romani RW (2007) Birth kick distributions and the spin-kick correlation of young pulsars. Astrophys J 660:1357-1374. https://doi.org/10.1086/513597
- Nishikawa H, Kitamura K (2008) Very simple, carbuncle-free, boundary-layer-resolving, rotated-hybrid Riemann solvers. J Comput Phys 227(4):2560-2581. https://doi.org/10.1016/j.jcp.2007.11.003
- Nomoto K, Leung SC (2017) Thermonuclear explosions of Chandrasekhar mass white dwarfs. In: Alsabti A, Murdin P (eds) Handbook of supernovae. Springer, Cham, p 1275-1330. https://doi.org/10.1007/ 978-3-319-21846-5\_62

<!-- image -->

- Nomoto K, Tominaga N, Umeda H, Kobayashi C, Maeda K (2006) Nucleosynthesis yields of corecollapse supernovae and hypernovae, and galactic chemical evolution. Nucl Phys A 777:424-458. https://doi.org/10.1016/j.nuclphysa.2006.05.008
- Nonaka A, Almgren AS, Bell JB, Lijewski MJ, Malone CM, Zingale M (2010) MAESTRO: an adaptive low mach number hydrodynamics algorithm for stellar flows. Astrophys J Suppl 188(2):358-383. https://doi.org/10.1088/0067-0049/188/2/358
- Nonaka A, Aspden AJ, Zingale M, Almgren AS, Bell JB, Woosley SE (2012) High-resolution simulations of convection preceding ignition in type Ia supernovae using adaptive mesh refinement. Astrophys J 745(1):73. https://doi.org/10.1088/0004-637X/745/1/73
- Nordhaus J, Brandt TD, Burrows A, Livne E, Ott CD (2010a) Theoretical support for the hydrodynamic mechanism of pulsar kicks. Phys Rev D 82(10):103016. https://doi.org/10.1103/PhysRevD.82. 103016
- Nordhaus J, Burrows A, Almgren A, Bell J (2010b) Dimension as a key to the neutrino mechanism of core-collapse supernova explosions. Astrophys J 720:694-703. https://doi.org/10.1088/0004-637X/ 720/1/694. arXiv:1006.3792
- Nordhaus J, Brandt TD, Burrows A, Almgren A (2012) The hydrodynamic origin of neutron star kicks. Mon Not R Astron Soc 423:1805-1812. https://doi.org/10.1111/j.1365-2966.2012.21002.x
- Noutsos A, Schnitzeler DHFM, Keane EF, Kramer M, Johnston S (2013) Pulsar spin-velocity alignment: kinematic ages, birth periods and braking indices. Mon Not R Astron Soc 430:2281-2301. https:// doi.org/10.1093/mnras/stt047
- Obergaulinger M, Aloy MA ´ (2017) Protomagnetar and black hole formation in high-mass stars. Mon Not R Astron Soc 469(1):L43-L47. https://doi.org/10.1093/mnrasl/slx046
- Obergaulinger M, Aloy MA ´ (2020) Magnetorotational core collapse of possible GRB progenitors: I. Explosion mechanisms. Mon Not R Astron Soc 492(4):4613-4634. https://doi.org/10.1093/mnras/ staa096. arXiv:1909.01105
- Obergaulinger M, Aloy MA, Mu ¨ller E (2006) Axisymmetric simulations of magneto-rotational core collapse: dynamics and gravitational wave signcal. Astron Astrophys 450:1107-1134. https://doi. org/10.1051/0004-6361:20054306
- Obergaulinger M, Janka HT, Aloy MA (2014) Magnetic field amplification and magnetically supported explosions of collapsing, non-rotating stellar cores. Mon Not R Astron Soc 445:3169-3199. https:// doi.org/10.1093/mnras/stu1969
- O'Connor EP, Couch SM (2018a) Exploring fundamentally three-dimensional phenomena in high-fidelity simulations of core-collapse supernovae. Astrophys J 865:81. https://doi.org/10.3847/1538-4357/ aadcf7
- O'Connor EP, Couch SM (2018b) Two-dimensional core-collapse supernova explosions aided by general relativity with multidimensional neutrino transport. Astrophys J 854(1):63. https://doi.org/10.3847/ 1538-4357/aaa893
- Oertel M, Hempel M, Kla ¨hn T, Typel S (2017) Equations of state for supernovae and compact stars. Rev Mod Phys 89(1):015007. https://doi.org/10.1103/RevModPhys.89.015007
- Ohnishi N, Kotake K, Yamada S (2006) Numerical analysis of standing accretion shock instability with neutrino heating in supernova cores. Astrophys J 641:1018-1028. https://doi.org/10.1086/500554. arXiv:astro-ph/0509765
- Ono M, Nagataki S, Ito H, Lee SH, Mao J, Ma H, Tolstov A (2013) Matter mixing in aspherical corecollapse supernovae: a search for possible conditions for conveying 56 Ni into high velocity regions. Astrophys J 773(2):161. https://doi.org/10.1088/0004-637X/773/2/161
- Ott CD (2009) Topical review: the gravitational-wave signature of core-collapse supernovae. Class Quantum Grav 26(6):063001. https://doi.org/10.1088/0264-9381/26/6/063001. arXiv:0809.0695
- Ott CD, Dimmelmeier H, Marek A, Janka HT, Hawke I, Zink B, Schnetter E (2007) 3D collapse of rotating stellar iron cores in general relativity including deleptonization and a nuclear equation of state. Phys Rev Lett 98(26):261101. https://doi.org/10.1103/PhysRevLett.98.261101
- Ott CD, Abdikamalov E, O'Connor E, Reisswig C, Haas R, Kalmus P, Drasco S, Burrows A, Schnetter E (2012) Correlated gravitational wave and neutrino signals from general-relativistic rapidly rotating iron core collapse. Phys Rev D 86(2):024026. https://doi.org/10.1103/PhysRevD.86.024026. arXiv: 1204.0512
- Ott CD, Roberts LF, da Silva SA, Fedrow JM, Haas R, Schnetter E (2018) The progenitor dependence of core-collapse supernovae from three-dimensional simulations with progenitor models of 12-40 M /C12 . Astrophys J 855:L3. https://doi.org/10.3847/2041-8213/aaa967

<!-- image -->

- O ¨ zel F, Freire P (2016) Masses, radii, and the equation of state of neutron stars. Annu Rev Astron Astrophys 54:401-440. https://doi.org/10.1146/annurev-astro-081915-023322
- Pan KC, Liebendo ¨rfer M, Couch SM, Thielemann FK (2018) Equation of state dependent dynamics and multi-messenger signals from stellar-mass black hole formation. Astrophys J 857(1):13. https://doi. org/10.3847/1538-4357/aab71d
- Papish O, Nordhaus J, Soker N (2015) A call for a paradigm shift from neutrino-driven to jet-driven corecollapse supernova mechanisms. Mon Not R Astron Soc 448(3):2362-2367. https://doi.org/10.1093/ mnras/stv131
- Patat F (2017) Introduction to supernova polarimetry. In: Alsabti A, Murdin P (eds) Handbook of supernovae. Springer, Cham, pp 1275-1330. https://doi.org/10.1007/978-3-319-21846-5\_110
- Paxton B, Bildsten L, Dotter A, Herwig F, Lesaffre P, Timmes F (2011) Modules for experiments in stellar astrophysics (MESA). Astrophys J Suppl 192:3. https://doi.org/10.1088/0067-0049/192/1/3
- Paxton B, Schwab J, Bauer EB, Bildsten L, Blinnikov S, Duffell P, Farmer R, Goldberg JA, Marchant P, Sorokina E, Thoul A, Townsend RHD, Timmes FX (2018) Modules for experiments in stellar astrophysics (MESA): convective boundaries, element diffusion, and massive star explosions. Astrophys J Suppl 234:34. https://doi.org/10.3847/1538-4365/aaa5a8
- Pejcha O, Prieto JL (2015) On the intrinsic diversity of type II-plateau supernovae. Astrophys J 806:225. https://doi.org/10.1088/0004-637X/806/2/225
- Pejcha O, Thompson TA (2012) The physics of the neutrino mechanism of core-collapse supernovae. Astrophys J 746:106. https://doi.org/10.1088/0004-637X/746/1/106
- Pejcha O, Thompson TA (2015) The landscape of the neutrino mechanism of core-collapse supernovae: neutron star and black hole mass functions, explosion energies, and nickel yields. Astrophys J 801:90. https://doi.org/10.1088/0004-637X/801/2/90
- Peng X, Xiao F, Takahashi K (2006) Conservative constraint for a quasi-uniform overset grid on the sphere. Quart J R Meteorol Soc 132(616):979-996. https://doi.org/10.1256/qj.05.18
- Perna R, Soria R, Pooley D, Stella L (2008) How rapidly do neutron stars spin at birth? Constraints from archival X-ray observations of extragalactic supernovae. Mon Not R Astron Soc 384:1638-1648. https://doi.org/10.1111/j.1365-2966.2007.12821.x
- Plewa T, Mu ¨ller E (1999) The consistent multi-fluid advection method. Astron Astrophys 342:179-191 arXiv:astro-ph/9807241
- Pons JA, Reddy S, Prakash M, Lattimer JM, Miralles JA (1999) Evolution of proto-neutron stars. Astrophys J 513:780-804. https://doi.org/10.1086/306889. arXiv:astro-ph/9807040
- Popov SB, Turolla R (2012) Initial spin periods of neutron stars in supernova remnants. Astrophys Space Sci 341:457-464. https://doi.org/10.1007/s10509-012-1100-z
- Powell J, Mu ¨ller B (2019) Gravitational wave emission from 3D explosion models of core-collapse supernovae with low and normal explosion energies. Mon Not R Astron Soc 487(1):1178-1190. https://doi.org/10.1093/mnras/stz1304
- Price DJ (2012) Smoothed particle hydrodynamics and magnetohydrodynamics. J Comput Phys 231(3):759-794. https://doi.org/10.1016/j.jcp.2010.12.011
- Proctor MRE (1981) Steady subcritical thermohaline convection. J Fluid Mech 105:507-521. https://doi. org/10.1017/S0022112081003315
- Quirk JJ (1994) A contribution to the great Riemann solver debate. Int J Num Meth Fluids 18:555-574. https://doi.org/10.1002/fld.1650180603
- Radice D, Couch SM, Ott CD (2015) Implicit large eddy simulations of anisotropic weakly compressible turbulence with application to core-collapse supernovae. Comput Astrophys Cosmol 2:7. https://doi. org/10.1186/s40668-015-0011-0
- Radice D, Ott CD, Abdikamalov E, Couch SM, Haas R, Schnetter E (2016) Neutrino-driven convection in core-collapse supernovae: high-resolution simulations. Astrophys J 820:76. https://doi.org/10. 3847/0004-637X/820/1/76
- Radice D, Morozova V, Burrows A, Vartanyan D, Nagakura H (2019) Characterizing the gravitational wave signal from core-collapse supernovae. Astrophys J 876(1):L9. https://doi.org/10.3847/20418213/ab191a
- Rampp M, Janka HT (2000) Spherically symmetric simulation with Boltzmann neutrino transport of core collapse and postbounce evolution of a 15 M /C12 star. Astrophys J 539:L33-L36. https://doi.org/10. 1086/312837
- Rampp M, Janka HT (2002) Radiation hydrodynamics with neutrinos. Variable Eddington factor method for core-collapse supernova simulations. Astron Astrophys 396:361-392. https://doi.org/10.1051/ 0004-6361:20021398

<!-- image -->

- Rantsiou E, Burrows A, Nordhaus J, Almgren A (2011) Induced rotation in 3D simulations of core collapse supernovae: implications for pulsar spins. Astrophys J 732:57. https://doi.org/10.1088/ 0004-637X/732/1/57
- Reinecke M, Hillebrandt W, Niemeyer JC (2002) Refined numerical models for multidimensional type Ia supernova simulations. Astron Astrophys 386:936-943. https://doi.org/10.1051/0004-6361: 20020323
- Reisswig C, Haas R, Ott CD, Abdikamalov E, Mo ¨sta P, Pollney D, Schnetter E (2013) Three-dimensional general-relativistic hydrodynamic simulations of binary neutron star coalescence and stellar collapse with multipatch grids. Phys Rev D 87(6):064023. https://doi.org/10.1103/PhysRevD.87.064023
- Rembiasz T, Obergaulinger M, Cerda ´-Dura ´n P, Aloy MA ´ , Mu ¨ller E (2017) On the measurements of numerical viscosity and resistivity in Eulerian MHD codes. Astrophys J Suppl 230(2):18. https://doi. org/10.3847/1538-4365/aa6254
- Repetto S, Davies MB, Sigurdsson S (2012) Investigating stellar-mass black hole kicks. Mon Not R Astron Soc 425:2799-2809. https://doi.org/10.1111/j.1365-2966.2012.21549.x
- Richtmyer RD (1960) Taylor instability in shock acceleration of compressible fluids. Commun Pure Appl Math 13(2):297-319. https://doi.org/10.1002/cpa.3160130207
- Ritter C, Andrassy R, Co ˆte ´ B, Herwig F, Woodward PR, Pignatari M, Jones S (2018) Convective-reactive nucleosynthesis of K, Sc, Cl and p-process isotopes in O-C shell mergers. Mon Not R Astron Soc 474(1):L1-L6. https://doi.org/10.1093/mnrasl/slx126
- Roberts LF, Shen G, Cirigliano V, Pons JA, Reddy S, Woosley SE (2012) Protoneutron star cooling with convection: the effect of the symmetry energy. Phys Rev Lett 108(6):061103. https://doi.org/10. 1103/PhysRevLett.108.061103. arXiv:1112.0335
- Roberts LF, Ott CD, Haas R, O'Connor EP, Diener P, Schnetter E (2016) General relativistic threedimensional multi-group neutrino radiation-hydrodynamics simulations of core-collapse supernovae. Astrophys J 831:98. https://doi.org/10.3847/0004-637X/831/1/98
- Rodionov AV (2017) Artificial viscosity in Godunov-type schemes to cure the carbuncle phenomenon. J Comput Phys 345:308-329. https://doi.org/10.1016/j.jcp.2017.05.024
- Ronchi C, Iacono R, Paolucci PS (1996) The ''cubed sphere'': a new method for the solution of partial differential equations in spherical geometry. J Comput Phys 124(1):93-114. https://doi.org/10.1006/ jcph.1996.0047
- Rosswog S (2015) SPH methods in the modelling of compact objects. Living Rev Comput Astrophys 1:1. https://doi.org/10.1007/lrca-2015-1
- Sawai H, Kotake K, Yamada S (2005) Core-collapse supernovae with nonuniform magnetic fields. Astrophys J 631(1):446-455. https://doi.org/10.1086/432529. arXiv:astro-ph/0505611
- Schaal K, Bauer A, Chandrashekar P, Pakmor R, Klingenberg C, Springel V (2015) Astrophysical hydrodynamics with a high-order discontinuous Galerkin scheme and adaptive mesh refinement. Mon Not R Astron Soc 453(4):4278-4300. https://doi.org/10.1093/mnras/stv1859
- Scheck L, Plewa T, Janka HT, Kifonidis K, Mu ¨ller E (2004) Pulsar recoil by large-scale anisotropies in supernova explosions. Phys Rev Lett 92(1):011103. https://doi.org/10.1103/PhysRevLett.92.011103
- Scheck L, Kifonidis K, Janka HT, Mu ¨ller E (2006) Multidimensional supernova simulations with approximative neutrino transport. I. Neutron star kicks and the anisotropy of neutrino-driven explosions in two spatial dimensions. Astron Astrophys 457:963-986. https://doi.org/10.1051/00046361:20064855. arXiv:astro-ph/0601302
- Scheck L, Janka HT, Foglizzo T, Kifonidis K (2008) Multidimensional supernova simulations with approximative neutrino transport. II. Convection and the advective-acoustic cycle in the supernova core. Astron Astrophys 477:931-952. https://doi.org/10.1051/0004-6361:20077701. arXiv:0704. 3001
- Scheidegger S, Whitehouse SC, Ka ¨ppeli R, Liebendo ¨rfer M (2010) Gravitational waves from supernova matter. Class Quantum Grav 27(11):114101. https://doi.org/10.1088/0264-9381/27/11/114101. arXiv:0912.1455
- Schnetter E, Hawley SH, Hawke I (2004) Evolutions in 3D numerical relativity using fixed mesh refinement. Class Quantum Grav 21(6):1465-1488. https://doi.org/10.1088/0264-9381/21/6/014
- Sedov LI (1959) Similarity and dimensional methods in mechanics. Academic Press, New York
- Shibata M, Liu YT, Shapiro SL, Stephens BC (2006) Magnetorotational collapse of massive stellar cores to neutron stars: simulations in full general relativity. Phys Rev D 74(10):104026. https://doi.org/10. 1103/PhysRevD.74.104026. arXiv:astro-ph/0610840
- Shigeyama T, Nomoto K (1990) Theoretical light curve of SN 1987A and mixing of hydrogen and nickel in the ejecta. Astrophys J 360:242-256. https://doi.org/10.1086/169114

<!-- image -->

- Shimizu T, Yamada S, Sato K (1993) Three-dimensional simulations of convection in supernova cores. Publ Astron Soc Japan 45:L53-L57
- Shiota D, Kusano K, Miyoshi T, Shibata K (2010) Magnetohydrodynamic modeling for a formation process of coronal mass ejections: interaction between an ejecting flux rope and an ambient field. Astrophys J 718(2):1305-1314. https://doi.org/10.1088/0004-637X/718/2/1305
- Shu FH (1992) Physics of astrophysics, vol II. University Science Books, Mill Valley
- Shu CW (1997) Essentially non-oscillatory and weighted essentially non-oscillatory schemes for hyperbolic conservation laws. Technical report, Institute for Computer Applications in Science and Engineering (ICASE)
- Simon S, Mandal JC (2019) A simple cure for numerical shock instability in the HLLC Riemann solver. J Comput Phys 378:477-496. https://doi.org/10.1016/j.jcp.2018.11.022
- Skinner MA, Dolence JC, Burrows A, Radice D, Vartanyan D (2019) FORNAX: a flexible code for multiphysics astrophysical simulations. Astrophys J Suppl 241(1):7. https://doi.org/10.3847/15384365/ab007f
- Smartt SJ (2015) Observational constraints on the progenitors of core-collapse supernovae: the case for missing high-mass stars. Publ Astron Soc Australia 32:e016. https://doi.org/10.1017/pasa.2015.17
- Smartt SJ, Eldridge JJ, Crockett RM, Maund JR (2009) The death of massive stars: I. Observational constraints on the progenitors of type II-P supernovae. Mon Not R Astron Soc 395:1409-1437. https://doi.org/10.1111/j.1365-2966.2009.14506.x. arXiv:0809.0403
- Soker N (2019) Possible indications for jittering jets in core collapse supernova explosion simulations. arXiv:1907.13312
- Solanki SK, Inhester B, Schu ¨ssler M (2006) The solar magnetic field. Rep Progr Phys 69(3):563-668. https://doi.org/10.1088/0034-4885/69/3/R02
- Springel V (2010) E pur si muove: Galilean-invariant cosmological hydrodynamical simulations on a moving mesh. Mon Not R Astron Soc 401(2):791-851. https://doi.org/10.1111/j.1365-2966.2009. 15715.x
- Spruit HC (2013) Semiconvection: theory. Astron Astrophys 552:A76. https://doi.org/10.1051/00046361/201220575
- Spruit HC (2015) The growth of helium-burning cores. Astron Astrophys 582:L2. https://doi.org/10.1051/ 0004-6361/201527171
- Spruit H, Phinney ES (1998) Birth kicks as the origin of pulsar rotation. Nature 393:139-141. https://doi. org/10.1038/30168
- Stancliffe RJ, Dearborn DSP, Lattanzio JC, Heap SA, Campbell SW (2011) Three-dimensional hydrodynamical simulations of a proton ingestion episode in a low-metallicity asymptotic giant branch star. Astrophys J 742:121. https://doi.org/10.1088/0004-637X/742/2/121
- Staritsin EI (2013) Turbulent entrainment at the boundaries of the convective cores of main-sequence stars. Astron Rep 57:380-390. https://doi.org/10.1134/S1063772913050089
- Stone JM, Norman ML (1992) ZEUS-2D: a radiation magnetohydrodynamics code for astrophysical flows in two space dimensions. I. The hydrodynamic algorithms and tests. Astrophys J Suppl 80:753. https://doi.org/10.1086/191680
- Strang EJ, Fernando HJS (2001) Entrainment and mixing in stratified shear flows. J Fluid Mech 428:349-386
- Sukhbold T, Woosley SE (2014) The compactness of presupernova stellar cores. Astrophys J 783:10. https://doi.org/10.1088/0004-637X/783/1/10
- Sukhbold T, Ertl T, Woosley SE, Brown JM, Janka HT (2016) Core-collapse supernovae from 9 to 120 solar masses based on neutrino-powered explosions. Astrophys J 821:38. https://doi.org/10.3847/ 0004-637X/821/1/38
- Summa A, Hanke F, Janka HT, Melson T, Marek A, Mu ¨ller B (2016) Progenitor-dependent explosion dynamics in self-consistent, axisymmetric simulations of neutrino-driven core-collapse supernovae. Astrophys J 825:6. https://doi.org/10.3847/0004-637X/825/1/6
- Summa A, Janka HT, Melson T, Marek A (2018) Rotation-supported Neutrino-driven supernova explosions in three dimensions and the critical luminosity condition. Astrophys J 852:28. https://doi. org/10.3847/1538-4357/aa9ce8
- Suresh A, Huynh HT (1997) Accurate monotonicity-preserving schemes with Runge-Kutta time stepping. J Comput Phys 136(1):83-99. https://doi.org/10.1006/jcph.1997.5745
- Sutherland RS, Bisset DK, Bicknell GV (2003) The numerical simulation of radiative shocks. I. The elimination of numerical shock instabilities using a local oscillation filter. Astrophys J Suppl 147:187-195. https://doi.org/10.1086/374795

<!-- image -->

- Suwa Y, Kotake K, Takiwaki T, Whitehouse SC, Liebendo ¨rfer M, Sato K (2010) Explosion geometry of a rotating 13 M /C12 star driven by the SASI-aided neutrino-heating supernova mechanism. Publ Astron Soc Japan 62:L49 ? arXiv:0912.1157
- Suwa Y, Takiwaki T, Kotake K, Fischer T, Liebendo ¨rfer M, Sato K (2013) On the importance of the equation of state for the neutrino-driven supernova explosion mechanism. Astrophys J 764:99. https://doi.org/10.1088/0004-637X/764/1/99. arXiv:1206.6101
- Suwa Y, Yoshida T, Shibata M, Umeda H, Takahashi K (2015) Neutrino-driven explosions of ultrastripped type Ic supernovae generating binary neutron stars. Mon Not R Astron Soc 454:3073-3081. https://doi.org/10.1093/mnras/stv2195
- Suzuki TK, Sumiyoshi K, Yamada S (2008) Alfve ´n Wave-Driven Supernova Explosion. Astrophys J 678(2):1200-1206. https://doi.org/10.1086/533515
- Takahashi K, Yamada S (2014) Linear analysis on the growth of non-spherical perturbations in supersonic accretion flows. Astrophys J 794:162. https://doi.org/10.1088/0004-637X/794/2/162
- Takahashi K, Iwakami W, Yamamoto Y, Yamada S (2016) Links between the shock instability in corecollapse supernovae and asymmetric accretions of envelopes. Astrophys J 831(1):75. https://doi.org/ 10.3847/0004-637X/831/1/75
- Takiwaki T, Kotake K, Suwa Y (2012) Three-dimensional hydrodynamic core-collapse supernova simulations for an 11.2 M /C12 star with spectral neutrino transport. Astrophys J 749:98. https://doi.org/ 10.1088/0004-637X/749/2/98. arXiv:1108.3989
- Takiwaki T, Kotake K, Suwa Y (2014) A comparison of twoand three-dimensional neutrinohydrodynamics simulations of core-collapse supernovae. Astrophys J 786:83. https://doi.org/10. 1088/0004-637X/786/2/83
- Takiwaki T, Kotake K, Suwa Y (2016) Three-dimensional simulations of rapidly rotating core-collapse supernovae: finding a neutrino-powered explosion aided by non-axisymmetric flows. Mon Not R Astron Soc 461:L112-L116. https://doi.org/10.1093/mnrasl/slw105
- Tamborra I, Hanke F, Mu ¨ller B, Janka HT, Raffelt G (2013) Neutrino signature of supernova hydrodynamical instabilities in three dimensions. Phys Rev Lett 111(12):121104. https://doi.org/10. 1103/PhysRevLett.111.121104. arXiv:1307.7936
- Tamborra I, Hanke F, Janka HT, Mu ¨ller B, Raffelt GG, Marek A (2014a) Self-sustained asymmetry of lepton-number emission: a new phenomenon during the supernova shock-accretion phase in three dimensions. Astrophys J 792:96. https://doi.org/10.1088/0004-637X/792/2/96
- Tamborra I, Raffelt G, Hanke F, Janka HT, Mu ¨ller B (2014b) Neutrino emission characteristics and detection opportunities based on three-dimensional supernova simulations. Phys Rev D 90(4):045032. https://doi.org/10.1103/PhysRevD.90.045032
- Thompson C (2000) Accretional heating of asymmetric supernova cores. Astrophys J 534:915-933. https://doi.org/10.1086/308773
- Thompson TA, Quataert E, Burrows A (2005) Viscosity and rotation in core-collapse supernovae. Astrophys J 620:861-877. https://doi.org/10.1086/427177. arXiv:astro-ph/0403224
- Timmes FX (1999) Integration of nuclear reaction networks for stellar hydrodynamics. Astrophys J Suppl 124(1):241-263. https://doi.org/10.1086/313257
- Timmes FX, Hoffman RD, Woosley SE (2000) An inexpensive nuclear energy generation network for stellar hydrodynamics. Astrophys J Suppl 129(1):377-398. https://doi.org/10.1086/313407
- Toro EF (2009) Riemann solvers and numerical methods for fluid dynamics: a practical introduction, 3rd edn. Springer, Berlin. https://doi.org/10.1007/b79761
- Toro EF, Spruce M, Speares W (1994) Restoration of the contact surface in the HLL-Riemann solver. Shock Waves 4:25-34. https://doi.org/10.1007/BF01414629
- Ugliano M, Janka HT, Marek A, Arcones A (2012) Progenitor-explosion connection and remnant birth masses for neutrino-driven supernovae of iron-core progenitors. Astrophys J 757:69. https://doi.org/ 10.1088/0004-637X/757/1/69
- Umeda H, Nomoto K (2003) First-generation black-hole-forming supernovae and the metal abundance pattern of a very iron-poor star. Nature 422(6934):871-873. https://doi.org/10.1038/nature01571
- Utrobin VP, Wongwathanarat A, Janka HT, Mu ¨ller E (2015) Supernova 1987A: neutrino-driven explosions in three dimensions and light curves. Astron Astrophys 581:A40. https://doi.org/10.1051/ 0004-6361/201425513
- van den Horn LJ, van Weert CG (1983) Transport properties of neutrinos in stellar collapse. I. Bulk viscosity of collapsing stellar cores. Astron Astrophys 125(1):93-100
- van den Horn LJ, van Weert CG (1984) Transport properties of neutrinos in stellar collapse. II. Shear viscosity, heat conduction, and diffusion. Astron Astrophys 136:74-80

<!-- image -->

- van Leer B (1977) Towards the ultimate conservative difference scheme. IV. A new approach to numerical convection. J Comput Phys 23:276. https://doi.org/10.1016/0021-9991(77)90095-X
- Vartanyan D, Burrows A, Radice D (2019a) Temporal and angular variations of 3D core-collapse supernova emissions and their physical correlations. Mon Not R Astron Soc 489(2):2227-2246. https://doi.org/10.1093/mnras/stz2307
- Vartanyan D, Burrows A, Radice D, Skinner MA, Dolence J (2019b) A successful 3D core-collapse supernova explosion model. Mon Not R Astron Soc 482:351-369. https://doi.org/10.1093/mnras/ sty2585
- Viallet M, Baraffe I, Walder R (2011) Towards a new generation of multi-dimensional stellar evolution models: development of an implicit hydrodynamic code. Astron Astrophys 531:A86. https://doi.org/ 10.1051/0004-6361/201016374
- Viallet M, Meakin C, Arnett D, Moca ´k M (2013) Turbulent convection in stellar interiors. III. Mean-field analysis and stratification effects. Astrophys J 769:1. https://doi.org/10.1088/0004-637X/769/1/1
- Viallet M, Meakin C, Prat V, Arnett D (2015) Toward a consistent use of overshooting parametrizations in 1D stellar evolution codes. Astron Astrophys 580:A61. https://doi.org/10.1051/0004-6361/ 201526294
- Vigna-Go ´mez A, Neijssel CJ, Stevenson S, Barrett JW, Belczynski K, Justham S, de Mink SE, Mu ¨ller B, Podsiadlowski P, Renzo M, Sze ´csi D, Mandel I (2018) On the formation history of Galactic double neutron stars. Mon Not R Astron Soc 481:4009-4029. https://doi.org/10.1093/mnras/sty2463
- von Groote J (2014) General relativistic multi dimensional simulations of electron capture supernovae. PhD thesis, Technische Universita ¨t Mu ¨nchen. http://mediatum.ub.tum.de/?id=1227385
- Von Neumann J, Richtmyer RD (1950) A method for the numerical calculation of hydrodynamic shocks. J Appl Phys 21(3):232-237. https://doi.org/10.1063/1.1699639
- Walder R, Burrows A, Ott CD, Livne E, Lichtenstadt I, Jarrah M (2005) Anisotropies in the neutrino fluxes and heating profiles in two-dimensional, time-dependent, multigroup radiation hydrodynamics simulations of rotating core-collapse supernovae. Astrophys J 626:317-332. https://doi.org/10. 1086/429816. arXiv:astro-ph/0412187
- Wanajo S, Janka HT, Mu ¨ller B (2011) Electron-capture supernovae as the origin of elements beyond iron. Astrophys J 726:L15. https://doi.org/10.1088/2041-8205/726/2/L15
- Wang L, Wheeler JC (2008) Spectropolarimetry of supernovae. Annu Rev Astron Astrophys 46:433-474. https://doi.org/10.1146/annurev.astro.46.060407.145139
- Weaver TA, Zimmerman GB, Woosley SE (1978) Presupernova evolution of massive stars. Astrophys J 225:1021-1029. https://doi.org/10.1086/156569
- Weiss A, Hillebrandt W, Thomas HC, Ritter H (2004) Cox and Giuli's principles of stellar structure. Cambridge Scientific Publishers, Cambridge
- Wilson JR, Mayle RW (1988) Convection in core collapse supernovae. Phys Rep 163:63-77. https://doi. org/10.1016/0370-1573(88)90036-1
- Wilson JR, Mayle RW (1993) Report on the progress of supernova research by the Livermore group. Phys Rep 227:97-111. https://doi.org/10.1016/0370-1573(93)90059-M
- Winkler KHA, Norman ML, Mihalas D (1984) Adaptive-mesh radiation hydrodynamics: I. The radiation transport equation in a completely adaptive coordinate system. J Quant Spectrosc Radiat Transf 31:473-489. https://doi.org/10.1016/0022-4073(84)90054-2
- Winteler C, Ka ¨ppeli R, Perego A, Arcones A, Vasset N, Nishimura N, Liebendo ¨rfer M, Thielemann FK (2012) Magnetorotationally driven supernovae as the origin of early galaxy r-process elements? Astrophys J 750:L22. https://doi.org/10.1088/2041-8205/750/1/L22
- Wongwathanarat A (2019) A generalized solution for parallelized computation of the three-dimensional gravitational potential on a multipatch grid in spherical geometry. Astrophys J 875(2):118. https:// doi.org/10.3847/1538-4357/ab1263
- Wongwathanarat A, Hammer NJ, Mu ¨ller E (2010a) An axis-free overset grid in spherical polar coordinates for simulating 3D self-gravitating flows. Astron Astrophys 514:A48. https://doi.org/10. 1051/0004-6361/200913435
- Wongwathanarat A, Janka HT, Mu ¨ller E (2010b) Hydrodynamical neutron star kicks in three dimensions. Astrophys J 725:L106-L110. https://doi.org/10.1088/2041-8205/725/1/L106. arXiv:1010.0167
- Wongwathanarat A, Janka HT, Mu ¨ller E (2013) Three-dimensional neutrino-driven supernovae: neutron star kicks, spins, and asymmetric ejection of nucleosynthesis products. Astron Astrophys 552:A126. https://doi.org/10.1051/0004-6361/201220636

<!-- image -->

- Wongwathanarat A, Mu ¨ller E, Janka HT (2015) Three-dimensional simulations of core-collapse supernovae: from shock revival to shock breakout. Astron Astrophys 577:A48. https://doi.org/10. 1051/0004-6361/201425025. arXiv:1409.5431
- Wongwathanarat A, Grimm-Strele H, Mu ¨ller E (2016) APSARA: a multi-dimensional unsplit fourthorder explicit Eulerian hydrodynamics code for arbitrary curvilinear grids. Astron Astrophys 595:A41. https://doi.org/10.1051/0004-6361/201628205
- Wongwathanarat A, Janka HT, Mu ¨ller E, Pllumbi E, Wanajo S (2017) Production and distribution of 44 Ti and 56 Ni in a three-dimensional supernova model resembling cassiopeia A. Astrophys J 842(1):13. https://doi.org/10.3847/1538-4357/aa72de
- Woodward PR, Porter D, Dai W, Fuchsa T, Nowatzkia T, Knox M, Dimonte G, Falk Herwig F, Fryer C (2010) the piecewise-parabolic Boltzmann advection scheme (PPB) Applied to multifluid hydrodynamics. Los Alamos Natl. Lab. report LA-UR 10-01823. http://www.lcse.umn.edu/ PPMplusPPB
- Woodward PR, Herwig F, Lin PH (2014) Hydrodynamic simulations of h entrainment at the top of heshell flash convection. Astrophys J 798(1):49. https://doi.org/10.1088/0004-637x/798/1/49
- Woodward PR, Herwig F, Wetherbee T (2018) Simulating stellar hydrodynamics at extreme scale. Comput Sci Eng 20(5):8-17. https://doi.org/10.1109/MCSE.2018.05329811
- Woodward PR, Lin PH, Mao H, Andrassy R, Herwig F (2019) Simulating 3-D stellar hydrodynamics using PPM and PPB multifluid gas dynamics on CPU and CPU ? GPU nodes. J Phys: Conf Ser 1225:012020. https://doi.org/10.1088/1742-6596/1225/1/012020. arXiv:1810.13416
- Woosley SE, Bloom JS (2006) The supernova gamma-ray burst connection. Annu Rev Astron Astrophys 44(1):507-556. https://doi.org/10.1146/annurev.astro.43.072103.150558
- Woosley SE, Heger A (2006) The progenitor stars of gamma-ray bursts. Astrophys J 637:914-921. https://doi.org/10.1086/498500. arXiv:astro-ph/0508175
- Woosley SE, Heger A (2007) Nucleosynthesis and remnants in massive stars of solar metallicity. Phys Rep 442:269-283. https://doi.org/10.1016/j.physrep.2007.02.009. arXiv:astro-ph/0702176
- Woosley SE, Heger A (2012) Long gamma-ray transients from collapsars. Astrophys J 752:32. https://doi. org/10.1088/0004-637X/752/1/32
- Woosley SE, Heger A (2015) The remarkable deaths of 9-11 solar mass stars. Astrophys J 810:34. https:// doi.org/10.1088/0004-637X/810/1/34
- Woosley SE, Arnett WD, Clayton DD (1972) Hydrostatic oxygen burning in stars. II. Oxygen burning at balanced power. Astrophys J 175:731. https://doi.org/10.1086/151594
- Woosley SE, Arnett WD, Clayton DD (1973) The explosive burning of oxygen and silicon. Astrophys J Suppl 26:231. https://doi.org/10.1086/190282
- Yabe T, Hoshino H, Tsuchiya T (1991) Two- and three-dimensional behavior of Rayleigh-Taylor and Kelvin-Helmholtz instabilities. Phys Rev A 44:2756-2758. https://doi.org/10.1103/PhysRevA.44. 2756
- Yadav N, Mu ¨ller B, Janka HT, Melson T, Heger A (2020) Large-scale mixing in a violent oxygen-neon shell merger prior to a core-collapse supernova. Astrophys J 890(2):94. https://doi.org/10.3847/ 1538-4357/ab66bb. arXiv:1905.04378
- Yakunin KN, Marronetti P, Mezzacappa A, Bruenn SW, Lee C, Chertkow MA, Hix WR, Blondin JM, Lentz EJ, Bronson Messer OE, Yoshida S (2010) Gravitational waves from core collapse supernovae. Class Quantum Grav 27(19):194005. https://doi.org/10.1088/0264-9381/27/19/194005. arXiv:1005.0779
- Yamada S, Sawai H (2004) Numerical study on the rotational collapse of strongly magnetized cores of massive stars. Astrophys J 608(2):907-924. https://doi.org/10.1086/420760
- Yamada S, Shimizu T, Sato K (1993) Convective instability in hot bubble in a delayed supernova explosion. Progr Theor Phys 89:1175-1182. https://doi.org/10.1143/PTP.89.1175
- Yamada S, Janka HT, Suzuki H (1999) Neutrino transport in type II supernovae: Boltzmann solver vs. Monte Carlo method. Astron Astrophys 344:533-550 arXiv:astro-ph/9809009
- Yamasaki T, Foglizzo T (2008) Effect of rotation on the stability of a stalled cylindrical shock and its consequences for core-collapse supernovae. Astrophys J 679(1):607-615. https://doi.org/10.1086/ 587732
- Yamasaki T, Yamada S (2007) Stability of accretion flows with stalled shocks in core-collapse supernovae. Astrophys J 656:1019-1037. https://doi.org/10.1086/510505. arXiv:astro-ph/0606581

<!-- image -->

- Yoon SC, Chun W, Tolstov A, Blinnikov S, Dessart L (2019) Type Ib/Ic supernovae: effect of nickel mixing on the early-time color evolution and implications for the progenitors. Astrophys J 872(2):174. https://doi.org/10.3847/1538-4357/ab0020
- Yoshida T, Takiwaki T, Kotake K, Takahashi K, Nakamura K, Umeda H (2019) One-, two-, and threedimensional simulations of oxygen-shell burning just before the core collapse of massive stars. Astrophys J 881(1):16. https://doi.org/10.3847/1538-4357/ab2b9d
- Young PA, Arnett D (2005) Observational tests and predictive stellar evolution. II. Nonstandard models. Astrophys J 618(2):908-918. https://doi.org/10.1086/426131
- Young PA, Meakin C, Arnett D, Fryer CL (2005) The impact of hydrodynamic mixing on supernova progenitors. Astrophys J 629(2):L101-L104. https://doi.org/10.1086/447769
- Yudin AV, Nadyozhin DK (2008) The approximation of neutrino heat conduction with neutrino scattering. Astron Lett 34(3):198-209. https://doi.org/10.1007/s11443-008-3007-0
- Zahn JP (1992) Circulation and turbulence in rotating stars. Astron Astrophys 265:115-132
- Zhou Y (2017) Rayleigh-Taylor and Richtmyer-Meshkov instability induced flow, turbulence, and mixing. I. Phys Rep 720:1-136. https://doi.org/10.1016/j.physrep.2017.07.005
- Zingale M, Dursi LJ, ZuHone J, Calder AC, Fryxell B, Plewa T, Truran JW, Caceres A, Olson K, Ricker PM, Riley K, Rosner R, Siegel A, Timmes FX, Vladimirova N (2002) Mapping Initial hydrostatic models in Godunov codes. Astrophys J Suppl 143(2):539-565. https://doi.org/10.1086/342754
- Zingale M, Nonaka A, Almgren AS, Bell JB, Malone CM, Woosley SE (2011) The convective phase preceding type Ia supernovae. Astrophys J 740(1):8. https://doi.org/10.1088/0004-637X/740/1/8

Publisher's Note Springer Nature remains neutral with regard to jurisdictional claims in published maps and institutional affiliations.

<!-- image -->