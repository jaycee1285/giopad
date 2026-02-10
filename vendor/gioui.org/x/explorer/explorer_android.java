package org.gioui.x.explorer;

import android.content.Context;
import android.util.Log;
import android.content.Intent;
import android.view.View;
import android.app.Activity;
import android.Manifest;
import android.content.pm.PackageManager;
import android.os.Handler.Callback;
import android.os.Handler;
import android.net.Uri;
import android.app.Fragment;
import android.app.FragmentManager;
import android.app.FragmentTransaction;
import android.os.Looper;
import android.content.ContentResolver;
import java.io.InputStream;
import java.io.OutputStream;
import android.webkit.MimeTypeMap;
import android.provider.DocumentsContract;
import android.database.Cursor;
import java.io.ByteArrayOutputStream;
import java.io.Closeable;
import java.io.Flushable;
import java.util.ArrayList;
import java.util.List;

public class explorer_android {
    final Fragment frag = new explorer_android_fragment();

    // List of requestCode used in the callback, to identify the caller.
    static List<Integer> import_codes = new ArrayList<Integer>();
    static List<Integer> export_codes = new ArrayList<Integer>();
    static List<Integer> dir_codes = new ArrayList<Integer>();

    // Functions defined on Golang.
    static public native void ImportCallback(InputStream f, int id, String err);
    static public native void ExportCallback(OutputStream f, int id, String err);
    static public native void DirCallback(String uri, int id, String err);

    public static class explorer_android_fragment extends Fragment {
        Context context;

        @Override public void onAttach(Context ctx) {
            context = ctx;
            super.onAttach(ctx);
        }

        @Override public void onActivityResult(int requestCode, int resultCode, Intent data) {
            super.onActivityResult(requestCode, resultCode, data);

            Activity activity = this.getActivity();

            activity.runOnUiThread(new Runnable() {
                public void run() {
                    if (import_codes.contains(Integer.valueOf(requestCode))) {
                        import_codes.remove(Integer.valueOf(requestCode));
                        if (resultCode != Activity.RESULT_OK) {
                            explorer_android.ImportCallback(null, requestCode, "");
                            activity.getFragmentManager().popBackStack();
                            return;
                        }
                        try {
                            InputStream f = activity.getApplicationContext().getContentResolver().openInputStream(data.getData());
                            explorer_android.ImportCallback(f, requestCode, "");
                        } catch (Exception e) {
                            explorer_android.ImportCallback(null, requestCode, e.toString());
                            return;
                        }
                    }

                    if (export_codes.contains(Integer.valueOf(requestCode))) {
                        export_codes.remove(Integer.valueOf(requestCode));
                        if (resultCode != Activity.RESULT_OK) {
                            explorer_android.ExportCallback(null, requestCode, "");
                            activity.getFragmentManager().popBackStack();
                            return;
                        }
                        try {
                            OutputStream f = activity.getApplicationContext().getContentResolver().openOutputStream(data.getData(), "wt");
                            explorer_android.ExportCallback(f, requestCode, "");
                        } catch (Exception e) {
                            explorer_android.ExportCallback(null, requestCode, e.toString());
                            return;
                        }
                    }

                    if (dir_codes.contains(Integer.valueOf(requestCode))) {
                        dir_codes.remove(Integer.valueOf(requestCode));
                        if (resultCode != Activity.RESULT_OK || data == null) {
                            explorer_android.DirCallback(null, requestCode, "");
                            activity.getFragmentManager().popBackStack();
                            return;
                        }
                        try {
                            Uri treeUri = data.getData();
                            // Take persistable permissions for this tree
                            int takeFlags = Intent.FLAG_GRANT_READ_URI_PERMISSION | Intent.FLAG_GRANT_WRITE_URI_PERMISSION;
                            activity.getContentResolver().takePersistableUriPermission(treeUri, takeFlags);
                            explorer_android.DirCallback(treeUri.toString(), requestCode, "");
                        } catch (Exception e) {
                            explorer_android.DirCallback(null, requestCode, e.toString());
                            return;
                        }
                    }
                }
            });

        }
    }

    public void exportFile(View view, String ext, int id) {
        askPermission(view);

        ((Activity) view.getContext()).runOnUiThread(new Runnable() {
            public void run() {
                registerFrag(view);
                export_codes.add(Integer.valueOf(id));
                
                final Intent intent = new Intent(Intent.ACTION_CREATE_DOCUMENT);
                intent.setType(MimeTypeMap.getSingleton().getMimeTypeFromExtension(ext));
                intent.addCategory(Intent.CATEGORY_OPENABLE);
                frag.startActivityForResult(Intent.createChooser(intent, ""), id);
            }
        });
    }

    public void importFile(View view, String mime, int id) {
        askPermission(view);

        ((Activity) view.getContext()).runOnUiThread(new Runnable() {
            public void run() {
                registerFrag(view);
                import_codes.add(Integer.valueOf(id));

                final Intent intent = new Intent(Intent.ACTION_GET_CONTENT);
                intent.setType("*/*");
                intent.addCategory(Intent.CATEGORY_OPENABLE);

                if (mime != null) {
                    final String[] mimes = mime.split(",");
                    if (mimes != null && mimes.length > 0) {
                        intent.putExtra(Intent.EXTRA_MIME_TYPES, mimes);
                    }
                }
                frag.startActivityForResult(Intent.createChooser(intent, ""), id);
            }
        });
    }

    public void importDir(View view, int id) {
        askPermission(view);

        ((Activity) view.getContext()).runOnUiThread(new Runnable() {
            public void run() {
                registerFrag(view);
                dir_codes.add(Integer.valueOf(id));

                final Intent intent = new Intent(Intent.ACTION_OPEN_DOCUMENT_TREE);
                intent.addFlags(Intent.FLAG_GRANT_READ_URI_PERMISSION);
                intent.addFlags(Intent.FLAG_GRANT_WRITE_URI_PERMISSION);
                intent.addFlags(Intent.FLAG_GRANT_PERSISTABLE_URI_PERMISSION);
                frag.startActivityForResult(intent, id);
            }
        });
    }

    public void registerFrag(View view) {
        final Context ctx = view.getContext();
        final FragmentManager fm;

        try {
            fm = (FragmentManager) ctx.getClass().getMethod("getFragmentManager").invoke(ctx);
        } catch (Exception e) {
            e.printStackTrace();
            return;
        }

        if (fm.findFragmentByTag("explorer_android_fragment") != null) {
            return; // Already exists;
        }

        FragmentTransaction ft = fm.beginTransaction();
        ft.add(frag, "explorer_android_fragment");
        ft.commitNow();
    }

    public void askPermission(View view) {
        Activity activity = (Activity) view.getContext();

        if (activity.checkSelfPermission(Manifest.permission.READ_EXTERNAL_STORAGE) != PackageManager.PERMISSION_GRANTED) {
            activity.requestPermissions(new String[] { Manifest.permission.READ_EXTERNAL_STORAGE }, 255);
        }

        if (activity.checkSelfPermission(Manifest.permission.WRITE_EXTERNAL_STORAGE) != PackageManager.PERMISSION_GRANTED) {
            activity.requestPermissions(new String[] { Manifest.permission.WRITE_EXTERNAL_STORAGE }, 254);
        }
    }

    // List children of a document tree URI
    // Returns newline-separated entries: "type|name|uri" where type is "d" for dir, "f" for file
    public static String listDir(Context ctx, String treeUriStr) {
        try {
            Uri treeUri = Uri.parse(treeUriStr);
            Uri childrenUri = DocumentsContract.buildChildDocumentsUriUsingTree(
                treeUri, DocumentsContract.getTreeDocumentId(treeUri));

            ContentResolver resolver = ctx.getContentResolver();
            StringBuilder result = new StringBuilder();

            String[] projection = {
                DocumentsContract.Document.COLUMN_DOCUMENT_ID,
                DocumentsContract.Document.COLUMN_DISPLAY_NAME,
                DocumentsContract.Document.COLUMN_MIME_TYPE
            };

            Cursor cursor = resolver.query(childrenUri, projection, null, null, null);
            if (cursor != null) {
                while (cursor.moveToNext()) {
                    String docId = cursor.getString(0);
                    String name = cursor.getString(1);
                    String mimeType = cursor.getString(2);

                    boolean isDir = DocumentsContract.Document.MIME_TYPE_DIR.equals(mimeType);
                    Uri docUri = DocumentsContract.buildDocumentUriUsingTree(treeUri, docId);

                    if (result.length() > 0) result.append("\n");
                    result.append(isDir ? "d" : "f");
                    result.append("|");
                    result.append(name);
                    result.append("|");
                    result.append(docUri.toString());
                }
                cursor.close();
            }
            return result.toString();
        } catch (Exception e) {
            return "ERROR:" + e.toString();
        }
    }

    // List children of a subdirectory within a tree
    public static String listSubDir(Context ctx, String treeUriStr, String docUriStr) {
        try {
            Uri treeUri = Uri.parse(treeUriStr);
            Uri docUri = Uri.parse(docUriStr);
            String docId = DocumentsContract.getDocumentId(docUri);
            Uri childrenUri = DocumentsContract.buildChildDocumentsUriUsingTree(treeUri, docId);

            ContentResolver resolver = ctx.getContentResolver();
            StringBuilder result = new StringBuilder();

            String[] projection = {
                DocumentsContract.Document.COLUMN_DOCUMENT_ID,
                DocumentsContract.Document.COLUMN_DISPLAY_NAME,
                DocumentsContract.Document.COLUMN_MIME_TYPE
            };

            Cursor cursor = resolver.query(childrenUri, projection, null, null, null);
            if (cursor != null) {
                while (cursor.moveToNext()) {
                    String childDocId = cursor.getString(0);
                    String name = cursor.getString(1);
                    String mimeType = cursor.getString(2);

                    boolean isDir = DocumentsContract.Document.MIME_TYPE_DIR.equals(mimeType);
                    Uri childDocUri = DocumentsContract.buildDocumentUriUsingTree(treeUri, childDocId);

                    if (result.length() > 0) result.append("\n");
                    result.append(isDir ? "d" : "f");
                    result.append("|");
                    result.append(name);
                    result.append("|");
                    result.append(childDocUri.toString());
                }
                cursor.close();
            }
            return result.toString();
        } catch (Exception e) {
            return "ERROR:" + e.toString();
        }
    }

    // Read file contents from a document URI
    public static byte[] readFile(Context ctx, String docUriStr) {
        try {
            Uri docUri = Uri.parse(docUriStr);
            ContentResolver resolver = ctx.getContentResolver();
            InputStream is = resolver.openInputStream(docUri);
            if (is == null) return null;

            ByteArrayOutputStream buffer = new ByteArrayOutputStream();
            byte[] data = new byte[4096];
            int nRead;
            while ((nRead = is.read(data, 0, data.length)) != -1) {
                buffer.write(data, 0, nRead);
            }
            is.close();
            return buffer.toByteArray();
        } catch (Exception e) {
            return null;
        }
    }

    // Get display name for root of tree
    public static String getTreeName(Context ctx, String treeUriStr) {
        try {
            Uri treeUri = Uri.parse(treeUriStr);
            String docId = DocumentsContract.getTreeDocumentId(treeUri);
            Uri docUri = DocumentsContract.buildDocumentUriUsingTree(treeUri, docId);

            ContentResolver resolver = ctx.getContentResolver();
            String[] projection = { DocumentsContract.Document.COLUMN_DISPLAY_NAME };
            Cursor cursor = resolver.query(docUri, projection, null, null, null);
            if (cursor != null && cursor.moveToFirst()) {
                String name = cursor.getString(0);
                cursor.close();
                return name;
            }
            return "";
        } catch (Exception e) {
            return "";
        }
    }

    // Write file contents to a document URI
    public static boolean writeFile(Context ctx, String docUriStr, byte[] data) {
        try {
            Uri docUri = Uri.parse(docUriStr);
            ContentResolver resolver = ctx.getContentResolver();
            OutputStream os = resolver.openOutputStream(docUri, "wt");
            if (os == null) return false;

            os.write(data);
            os.close();
            return true;
        } catch (Exception e) {
            return false;
        }
    }
}