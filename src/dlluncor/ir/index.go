package ir

import(
  "dlluncor/ir/types"
)

type Index struct {

}

var allDocs = []*types.DocMetadata{
  {
    Title: "Google Play services",
    Id: "com.google.android.gms",
    Description: `
Google Play services is used to update Google apps and apps from Google Play.
This component provides core functionality like authentication to your Google services, synchronized contacts, access to all the latest user privacy settings, and higher quality, lower-powered location based services.
Google Play services also enhances your app experience. It speeds up offline searches, provides more immersive maps, and improves gaming experiences.
Apps may not work if you uninstall Google Play services.
    `,
  },
  {
    Title: "Temple Run",
    Id: "com.imangi.templerun",
    Description: `
The addictive mega-hit Temple Run is now out for Android! All your friends are playing it - can you beat their high scores?!
You've stolen the cursed idol from the temple, and now you have to run for your life to escape the Evil Demon Monkeys nipping at your heels. Test your reflexes as you race down ancient temple walls and along sheer cliffs. Swipe to turn, jump and slide to avoid obstacles, collect coins and buy power ups, unlock new characters, and see how far you can run!
"In every treasure hunting adventure movie there’s one scene in which the plucky hero finally gets his hands on the treasure but then has to navigate a maze of booby traps in order to get out alive. Temple Run is this scene and nothing else. And it’s amazing." - SlideToPlay.com

REVIEWS
★ "Most thrilling and fun running game in a while, possibly ever." - TheAppera.com
★ "A fast and frenzied experience." - IGN.com
★ "Very addicting… definitely a very different running game." - Appolicious.com
★ Voted by TouchArcade Forums as Game of the Week
★ One of TouchArcade's Best Games of the Month
★ Over 50 MILLION players worldwide!
    `,
  },
  {
    Title: "Twilight",
    Id: "com.urbandroid.lux",
    Description: `
Are you having troubles to fall asleep? Are your kids hyperactive when playing with the tablet before bed time?
Are you using your smart phone or tablet in the late evening? Twilight may be a solution for you!
Recent research suggests that exposure to blue light before sleep may distort your natural (circadian) rhythm and cause inability to fall asleep.
The cause is the photoreceptor in your eyes, called Melanopsin. This receptor is sensitive to a narrow band of blue light in the 460-480nm range which may suppress Melatonin production - a hormone responsible for your healthy sleep-wake cycles.
In experimental scientific studies it has been shown an average person reading on a tablet or smart phone for a couple of hours before bed time may find their sleep delayed by about an hour.
The Twilight app makes your device screen adapt to the time of the day. It filters the blue spectrum on your phone or tablet after sunset and protects your eyes with a soft and pleasant red filter. The filter intensity is smoothly adjusted to the sun cycle based on your local sunset and sunrise times.
Please read the basics on circadian rhythm and the role of melatonin:
http://en.wikipedia.org/wiki/Melatonin
http://en.wikipedia.org/wiki/Melanopsin
http://en.wikipedia.org/wiki/Circadian_rhythms
http://en.wikipedia.org/wiki/Circadian_rhythm_disorder
Permissions explained:
- location - to find out your current sunset/surise times
- running apps - to stop Twilight in selected apps
- write settings - to set back-light
- network - access smartlight (Philips HUE) to shield you household light from blue
Automation through Tasker or other:
https://sites.google.com/site/twilight4android/automation
Examples of related scientific research:
Amplitude Reduction and Phase Shifts of Melatonin, Cortisol and Other Circadian Rhythms after a Gradual Advance of Sleep and Light Exposure in Humans
Derk-Jan Dijk, Jeanne F. Duffy, Edward J. Silva, Theresa L. Shanahan, Diane B. Boivin, Charles A. Czeisler 2012
Exposure to Room Light before Bedtime Suppresses Melatonin Onset and Shortens Melatonin Duration in Humans
Joshua J. Gooley, Kyle Chamberlain, Kurt A. Smith, Sat Bir S. Khalsa, Shantha M. W. Rajaratnam, Eliza Van Reen, Jamie M. Zeitzer, Charles A. Czeisler, Steven W. 2011
Effect of Light on Human Circadian Physiology
Jeanne F. Duffy, Charles A. Czeisler 2009
Efficacy of a single sequence of intermittent bright light pulses for delaying circadian phase in humans
Claude Gronfier, Kenneth P. Wright, Richard E. Kronauer, Megan E. Jewett, Charles A. Czeisler 2009
Intrinsic period and light intensity determine the phase relationship between melatonin and sleep in humans
Kenneth P. Wright, Claude Gronfier, Jeanne F. Duffy, Charles A. Czeisler 2009
The Impact of Sleep Timing and Bright Light Exposure on Attentional Impairment during Night Work
Nayantara Santhi, Daniel Aeschbach, Todd S. Horowitz, Charles A. Czeisler 2008
Short-Wavelength Light Sensitivity of Circadian, Pupillary, and Visual Awareness in Humans Lacking an Outer Retina
Farhan H. Zaidi, Joseph T. Hull, Stuart N. Peirson, Katharina Wulff, Daniel Aeschbach, Joshua J. Gooley, George C. Brainard, Kevin Gregory-Evans, Joseph F. Rizzo, III, Charles A. Czeisler, Russell G. Foster, Merrick J. Moseley, Steven W. Lockley. 2007
High sensitivity of the human circadian melatonin rhythm to resetting by short wavelength light.
Lockley SW, Brainard GC, Czeisler CA. 2003
Sensitivity of the human circadian pacemaker to nocturnal light: melatonin phase resetting and suppression
Jamie M Zeitzer, Derk-Jan Dijk, Richard E Kronauer, Emery N Brown, Charles A Czeisler 2000
Search for more evidence in many scientific articles and empirical studies online. Keywords: Circadian, Melatonin, Melanopsin, Light, Wavelength...
Similar to f.lux on Windows or Red Shift on Linux.
`},
}

func (i *Index) Find(q *query) []*doc {
  docs := []*doc{}
  for _, docM := range allDocs {
    docs = append(docs, &doc{
      name: docM.Title,
      data: docM,
    })
  }
  return docs
}
