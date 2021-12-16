package hipstershop.copyright;

import hipstershop.AdService;
import hipstershop.Demo;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import java.util.LinkedList;
import java.util.List;

// A massive repository of very many photos.
public class StockPhotos {

    private final static long MAX_TIME_MS = 6500;
    private final static List<CopyrightPhoto> photos = createDatabase();
    private static final Logger logger = LogManager.getLogger(AdService.class);

    boolean isCopyright(Demo.Ad ad){
        boolean result = true;
        long start = System.currentTimeMillis();
        for (CopyrightPhoto photo : photos) {
            if(photo.matchesAd(ad)){
                result = false;
            }
            if(System.currentTimeMillis() - start > MAX_TIME_MS){
                logger.warn("Copyright check exceeded max threshold. Aborting.");
                break;
            }
        }
        return result;
    }

    private static List<CopyrightPhoto> createDatabase() {
        List<CopyrightPhoto> result = new LinkedList<>();
        for(int i=0; i < 5000; i++){
            result.add(new CopyrightPhoto("photo" + i));
        }
        return result;
    }

    static class CopyrightPhoto {

        private final String id;

        public CopyrightPhoto(String id) {
            this.id = id;
        }

        public boolean matchesAd(Demo.Ad ad) {
            boolean matches = false;
            for (CopyrightPhoto photo : photos) {
                if(photo.id.equals(ad.getRedirectUrl())){
                    matches = true;
                }
            }
            return matches;
        }
    }

}
