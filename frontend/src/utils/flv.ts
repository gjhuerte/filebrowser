/*
 * FLV playback for video.js.
 *
 * Browsers can't play Flash video natively, so FLV sources are remuxed
 * to fragmented MP4 in the browser with mpegts.js and fed to the normal
 * <video> element through Media Source Extensions. mpegts.js is loaded
 * on demand, so it only costs bandwidth when an FLV is opened.
 *
 * iOS Safari on iPhone has no MSE; there the player shows its normal
 * "can't play" message with a download button.
 */
import videojs from "video.js";
import type Mpegts from "mpegts.js";

export const FLV_TYPE = "video/x-flv";

let registered = false;
let errorMessage = "This video can't be played in your browser.";

const supportsMse = () =>
  typeof window !== "undefined" &&
  !!(window.MediaSource || (window as any).ManagedMediaSource);

interface Source {
  src: string;
  type?: string;
}

export function registerFlvSupport(message?: string) {
  if (message) errorMessage = message;
  if (registered) return;
  registered = true;

  const canPlay = (type: string) =>
    type.toLowerCase() === FLV_TYPE && supportsMse() ? "maybe" : "";

  const Html5 = videojs.getTech("Html5") as any;
  Html5.registerSourceHandler(
    {
      canPlayType: canPlay,
      canHandleSource: (source: Source) => canPlay(source.type ?? ""),
      handleSource(source: Source, tech: any) {
        let flv: Mpegts.Player | null = null;
        let disposed = false;

        import("mpegts.js").then(({ default: mpegts }) => {
          if (disposed) return;
          flv = mpegts.createPlayer(
            { type: "flv", url: source.src, isLive: false },
            {
              enableWorker: false,
              // Don't download the whole file up front; stay a bit ahead
              lazyLoad: true,
              lazyLoadMaxDuration: 3 * 60,
              seekType: "range",
            }
          );
          flv.on(mpegts.Events.ERROR, () => {
            const player = videojs.getPlayer(tech.options_?.playerId);
            player?.error({ code: 4, message: errorMessage } as any);
          });
          flv.attachMediaElement(tech.el());
          flv.load();
        });

        return {
          dispose() {
            disposed = true;
            if (!flv) return;
            flv.unload();
            flv.detachMediaElement();
            flv.destroy();
            flv = null;
          },
        };
      },
    },
    0
  );
}
